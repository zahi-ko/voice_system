package memory

import (
	"context"
	"path/filepath"
	"testing"

	"voice_system/internal/domain/audio"
)

func testAudio(id string) (audio.Audio, audio.AudioData) {
	item, err := audio.New(id, "test", audio.Metadata{
		Format:     audio.FormatWAV,
		SampleRate: 8000,
		BitDepth:   32,
		Duration:   1.0,
	})
	if err != nil {
		panic(err)
	}
	return item, audio.AudioData{0.1, -0.2, 0.3, -0.4}
}

func TestSaveGetRoundtrip(t *testing.T) {
	store, err := NewAudioStore()
	if err != nil {
		t.Fatal(err)
	}
	item, data := testAudio("id-1")
	ctx := context.Background()

	if err := store.Save(ctx, item, data); err != nil {
		t.Fatalf("save: %v", err)
	}

	gotItem, gotData, err := store.Get(ctx, "id-1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if gotItem.ID != item.ID || gotItem.Name != item.Name {
		t.Errorf("item = %+v, want %+v", gotItem, item)
	}
	if len(gotData) != len(data) {
		t.Fatalf("data length = %d, want %d", len(gotData), len(data))
	}
	for i := range data {
		if gotData[i] != data[i] {
			t.Errorf("sample %d = %f, want %f", i, gotData[i], data[i])
		}
	}
}

func TestGetInfo(t *testing.T) {
	store, _ := NewAudioStore()
	item, data := testAudio("id-1")
	ctx := context.Background()

	if err := store.Save(ctx, item, data); err != nil {
		t.Fatalf("save: %v", err)
	}

	info, err := store.GetInfo(ctx, "id-1")
	if err != nil {
		t.Fatalf("getInfo: %v", err)
	}
	if info.ID != "id-1" {
		t.Errorf("info id = %q", info.ID)
	}
}

func TestList(t *testing.T) {
	store, _ := NewAudioStore()
	ctx := context.Background()

	for _, id := range []string{"id-1", "id-2", "id-3"} {
		item, data := testAudio(id)
		if err := store.Save(ctx, item, data); err != nil {
			t.Fatalf("save %s: %v", id, err)
		}
	}

	list, err := store.List(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list.Items) != 3 {
		t.Fatalf("list length = %d, want 3", len(list.Items))
	}
}

func TestDelete(t *testing.T) {
	store, _ := NewAudioStore()
	item, data := testAudio("id-1")
	ctx := context.Background()

	if err := store.Save(ctx, item, data); err != nil {
		t.Fatalf("save: %v", err)
	}
	if err := store.Delete(ctx, "id-1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, _, err := store.Get(ctx, "id-1"); err == nil {
		t.Error("get after delete should fail")
	}
}

func TestGetMissing(t *testing.T) {
	store, _ := NewAudioStore()
	ctx := context.Background()

	if _, _, err := store.Get(ctx, "missing"); err == nil {
		t.Error("get missing should fail")
	}
	if _, err := store.GetInfo(ctx, "missing"); err == nil {
		t.Error("getInfo missing should fail")
	}
}

func TestReplace(t *testing.T) {
	store, _ := NewAudioStore()
	item, data := testAudio("id-1")
	ctx := context.Background()

	if err := store.Save(ctx, item, data); err != nil {
		t.Fatalf("save: %v", err)
	}

	replaced, _ := audio.New("id-1", "replaced", item.Meta)
	if err := store.Replace(ctx, "id-1", replaced, audio.AudioData{0.9, 0.8}); err != nil {
		t.Fatalf("replace: %v", err)
	}

	gotItem, gotData, err := store.Get(ctx, "id-1")
	if err != nil {
		t.Fatalf("get after replace: %v", err)
	}
	if gotItem.Name != "replaced" {
		t.Errorf("name after replace = %q, want %q", gotItem.Name, "replaced")
	}
	if len(gotData) != 2 || gotData[0] != 0.9 {
		t.Errorf("data after replace = %v, want [0.9 0.8]", gotData)
	}
}

func TestStorePersistsFilesUnderRoot(t *testing.T) {
	store, _ := NewAudioStore()
	item, data := testAudio("id-1")

	if err := store.Save(context.Background(), item, data); err != nil {
		t.Fatalf("save: %v", err)
	}
	// gob 序列化文件应存在于 root 目录下，文件名即音频 ID
	if _, err := lookupFile(store.root, "id-1"); err != nil {
		t.Errorf("audio file missing under root: %v", err)
	}
}

func lookupFile(root, id string) (string, error) {
	matches, err := filepath.Glob(filepath.Join(root, id))
	if err != nil || len(matches) == 0 {
		if err == nil {
			err = filepath.ErrBadPattern
		}
		return "", err
	}
	return matches[0], nil
}
