package postgres

import (
	"blog-gunk/blog/storage"
	"context"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPost(t *testing.T) {

	s := newTestStorage(t)

	tests := []struct {
		name    string
		in      storage.Post
		want    int64
		wantErr bool
	}{
		{
			name: "CREATE_POST_SUCCESS",
			in: storage.Post{
				CatID: 2,
				Title: "This is Title",
				Description: "This is description",
				Image: "This is image",
			},
			want: 1,
		},
		{
			name: "CREATE_POST_SUCCESS",
			in: storage.Post{
				CatID: 2,
				Title: "This is Title 2",
				Description: "This is description 2",
				Image: "This is image 2",
			},
			want: 2,
		},
		{
			name: "FAILED_DUPLICATE_TITLE",
			in: storage.Post{
				CatID: 2,
				Title: "This is Title",
				Description: "This is description",
				Image: "This is image",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.CreatePost(context.Background(), tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Storage.Create() = %v, want %v", got, tt.want)
			}

			
		})
	}
}

func TestListPost(t *testing.T) {
	s := newTestStorage(t)

	tests := []struct {
		name    string
		in      int64
		want    []storage.Post
		wantErr bool
	}{
		{
			name: "GET_ALL_POST_SUCCESS",
			want: []storage.Post{
				{
					ID: 1,
					CatID: 2,
					Title: "This is Title",
					Description: "This is description",
					Image: "This is image",
					CatName: "this is title 2",
				},{
					ID: 2,
					CatID: 2,
					Title: "This is Title 2",
					Description: "This is description 2",
					Image: "This is image 2",
					CatName: "this is title 2",
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			gotList, err := s.ListPost(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			sort.Slice(tt.want, func(i, j int) bool {
				return tt.want[i].ID < tt.want[j].ID
			})

			sort.Slice(gotList, func(i, j int) bool {
				return gotList[i].ID < gotList[j].ID
			})

			for i, got := range gotList {

				if !cmp.Equal(got, tt.want[i]) {
					t.Errorf("Diff: got -, want += %v", cmp.Diff(got, tt.want[i]))
				}

			}

		})
	}
}

func TestGetPost(t *testing.T) {

	s := newTestStorage(t)

	tests := []struct {
		name    string
		in      int64
		want    storage.Post
		wantErr bool
	}{
		{
			name: "GET_POST_SUCCESS",
			in: 1,
			want: storage.Post{
				ID: 1,
				CatID: 2,
				Title: "This is Title",
				Description: "This is description",
				Image: "This is image",
				CatName: "this is title 2",
			},
		},
		{
			name: "GET_POST_SUCCESS",
			in: 2,
			want: storage.Post{
				ID: 2,
				CatID: 2,
				Title: "This is Title 2",
				Description: "This is description 2",
				Image: "This is image 2",
				CatName: "this is title 2",
			},
		},
		{
			name: "FAILED_TO_GET_POST",
			in: 3,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetPost(context.Background(), tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("Diff: got -, want += %v", cmp.Diff(got, tt.want))
			}
			
		})
	}
}

func TestUpdatePost(t *testing.T) {

	s := newTestStorage(t)

	tests := []struct {
		name    string
		in      storage.Post
		wantErr bool
	}{
		{
			name: "UPDATE_POST_SUCCESS",
			in: storage.Post{
				ID: 1,
				CatID: 2,
				Title: "This is Title updated",
				Description: "This is description updated",
				Image: "This is image updated",
				
			},
		},
		{
			name: "FAILED_TO_UPDATE_POST",
			in: storage.Post{
				ID: 4,
				CatID: 1,
				Title: "This is Title 3",
				Description: "This is description 3",
				Image: "This is image 3",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := s.UpdatePost(context.Background(), tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDeletePost(t *testing.T) {

	s := newTestStorage(t)

	tests := []struct {
		name    string
		in      int64
		want    bool
		wantErr bool
	}{
		{
			name: "DELETE_POST_SUCCESS",
			in: 1,
			want: true,
		},
		{
			name: "FAILED_TO_DELETE_POST",
			in: 3,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := s.DeletePost(context.Background(), tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}