package postgres

import (
	"blog-gunk/blog/storage"
	"context"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCreateCategory(t *testing.T) {

	s := newTestStorage(t)

	tests := []struct {
		name    string
		in      storage.Category
		want    int64
		wantErr bool
	}{
		{
			name: "CREATE_CATAGORY_SUCCESS",
			in: storage.Category{
				Name: "this is title",
			},
			want: 1,
		},
		{
			name: "CREATE_CATAGORY_SUCCESS",
			in: storage.Category{
				Name: "this is title 2",
			},
			want: 2,
		},
		{
			name: "FAILED_DUPLICATE_TITLE",
			in: storage.Category{
				Name: "this is title",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Create(context.Background(), tt.in)
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

func TestGetCategory(t *testing.T) {

	s := newTestStorage(t)

	tests := []struct {
		name    string
		in      int64
		want    storage.Category
		wantErr bool
	}{
		{
			name: "GET_CATAGORY_SUCCESS",
			in: 1,
			want: storage.Category{
				ID: 1,
				Name: "this is title",
				Status: false,
			},
		},
		{
			name: "GET_CATAGORY_SUCCESS",
			in: 2,
			want: storage.Category{
				ID: 2,
				Name: "this is title 2",
				Status: false,
			},
		},
		{
			name: "FAILED_TO_GET_CATAGORY",
			in: 3,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Get(context.Background(), tt.in)
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

func TestGetsCategory(t *testing.T) {

	s := newTestStorage(t)

	tests := []struct {
		name    string
		in      int64
		want    []storage.Category
		wantErr bool
	}{
		{
			name: "GET_ALL_CATAGORY_SUCCESS",
			want: []storage.Category{
				{
					ID: 1,
					Name: "this is title",
					Status: false,
				},{
					ID: 2,
					Name: "this is title 2",
					Status: false,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			gotList, err := s.Gets(context.Background())
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

func TestUpdateCategory(t *testing.T) {

	s := newTestStorage(t)

	tests := []struct {
		name    string
		in      storage.Category
		want    storage.Category
		wantErr bool
	}{
		{
			name: "UPDATE_CATAGORY_SUCCESS",
			in: storage.Category{
				ID: 1,
				Name: "this is title updated",
				Status: false,
			},
			want: storage.Category{
				ID: 1,
				Name: "this is title updated",
				Status: false,
			},
		},
		{
			name: "FAILED_TO_UPDATE_CATAGORY",
			in: storage.Category{
				ID: 3,
				Name: "this is title updated",
				Status: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := s.Update(context.Background(), tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDeleteCategory(t *testing.T) {

	s := newTestStorage(t)

	tests := []struct {
		name    string
		in      int64
		want    bool
		wantErr bool
	}{
		{
			name: "DELETE_CATAGORY_SUCCESS",
			in: 1,
			want: true,
		},
		{
			name: "FAILED_TO_DELETE_CATAGORY",
			in: 3,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := s.Delete(context.Background(), tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}