package bucket

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"analytics/internal/model"
)

func TestBucket_BuildIndex(t *testing.T) {
	type fields struct {
		data      []model.Event
		indexList map[string][]*model.Event
	}
	type args struct {
		fields []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string][]*model.Event
		wantErr bool
	}{
		{
			name: "test 1",
			fields: fields{
				data: []model.Event{
					{ID: 11185376329, Type: "PushEvent", ActorID: 53201765, RepoID: 224252202},
					{ID: 11185376330, Type: "PushEvent", ActorID: 8422699, RepoID: 224252202},
					{ID: 11185376331, Type: "PullRequestEvent", ActorID: 8422699, RepoID: 224252202},
					{ID: 11185376332, Type: "CreateEvent", ActorID: 53201765, RepoID: 231161852},
					{ID: 11185376333, Type: "PushEvent", ActorID: 8422600, RepoID: 224252200},
					{ID: 11185376334, Type: "CreateEvent", ActorID: 53201765, RepoID: 231161852},
				},
			},
			args: args{fields: []string{"Type", "RepoID"}},
			want: map[string][]*model.Event{
				"Type:PushEvent": []*model.Event{
					{ID: 11185376329, Type: "PushEvent", ActorID: 53201765, RepoID: 224252202},
					{ID: 11185376330, Type: "PushEvent", ActorID: 8422699, RepoID: 224252202},
					{ID: 11185376333, Type: "PushEvent", ActorID: 8422600, RepoID: 224252200},
				},
				"Type:PullRequestEvent": []*model.Event{
					{ID: 11185376331, Type: "PullRequestEvent", ActorID: 8422699, RepoID: 224252202},
				},
				"Type:CreateEvent": []*model.Event{
					{ID: 11185376332, Type: "CreateEvent", ActorID: 53201765, RepoID: 231161852},
					{ID: 11185376334, Type: "CreateEvent", ActorID: 53201765, RepoID: 231161852},
				},
				"RepoID:224252202": []*model.Event{
					{ID: 11185376329, Type: "PushEvent", ActorID: 53201765, RepoID: 224252202},
					{ID: 11185376330, Type: "PushEvent", ActorID: 8422699, RepoID: 224252202},
					{ID: 11185376331, Type: "PullRequestEvent", ActorID: 8422699, RepoID: 224252202},
				},
				"RepoID:231161852": []*model.Event{
					{ID: 11185376332, Type: "CreateEvent", ActorID: 53201765, RepoID: 231161852},
					{ID: 11185376334, Type: "CreateEvent", ActorID: 53201765, RepoID: 231161852},
				},
				"RepoID:224252200": []*model.Event{
					{ID: 11185376333, Type: "PushEvent", ActorID: 8422600, RepoID: 224252200},
				},
			},
		},
	}
	assert := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &EventBucket{
				data:      tt.fields.data,
				indexList: tt.fields.indexList,
			}
			err := s.BuildIndex(tt.args.fields...)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(tt.want, s.indexList)
		})
	}
}

func TestBucket_in(t *testing.T) {
	type args struct {
		in  []interface{}
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test 1",
			args: args{
				val: "PushEvent",
				in:  []interface{}{"PushEvent"},
			},
			want: true,
		},
		{
			name: "test 2",
			args: args{
				val: "ForkEvent",
				in:  []interface{}{"PushEvent"},
			},
			want: false,
		},
		{
			name: "test 3",
			args: args{
				val: "PushEvent",
				in:  []interface{}{"PullRequestEvent", "PushEvent"},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValueIn(tt.args.val, tt.args.in); got != tt.want {
				t.Errorf("isValueIn() = %v, want %v", got, tt.want)
			}
		})
	}
}
