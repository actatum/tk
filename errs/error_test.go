package errs

import (
	"fmt"
	"testing"
)

func TestError_Error(t *testing.T) {
	type fields struct {
		Kind    Kind
		Message string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "internal",
			fields: fields{
				Kind:    Internal,
				Message: "test",
			},
			want: "kind=internal message=test",
		},
		{
			name: "already exists",
			fields: fields{
				Kind:    AlreadyExists,
				Message: "test",
			},
			want: "kind=already_exists message=test",
		},
		{
			name: "not found",
			fields: fields{
				Kind:    NotFound,
				Message: "test",
			},
			want: "kind=not_found message=test",
		},
		{
			name: "unknown kind",
			fields: fields{
				Kind:    Kind(128),
				Message: "test",
			},
			want: "kind=internal message=test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				Kind:    tt.fields.Kind,
				Message: tt.fields.Message,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorKind(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want Kind
	}{
		{
			name: "internal",
			args: args{
				err: Errorf(Internal, "test"),
			},
			want: Internal,
		},
		{
			name: "already exists",
			args: args{
				err: Errorf(AlreadyExists, "test"),
			},
			want: AlreadyExists,
		},
		{
			name: "not found",
			args: args{
				err: Errorf(NotFound, "test"),
			},
			want: NotFound,
		},
		{
			name: "non package error",
			args: args{
				fmt.Errorf("error"),
			},
			want: Internal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrorKind(tt.args.err); got != tt.want {
				t.Errorf("ErrorKind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorMessage(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "internal",
			args: args{
				err: Errorf(Internal, "test"),
			},
			want: "test",
		},
		{
			name: "already exists",
			args: args{
				err: Errorf(AlreadyExists, "test"),
			},
			want: "test",
		},
		{
			name: "not found",
			args: args{
				err: Errorf(NotFound, "test"),
			},
			want: "test",
		},
		{
			name: "non package error",
			args: args{
				fmt.Errorf("error"),
			},
			want: "internal error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrorMessage(tt.args.err); got != tt.want {
				t.Errorf("ErrorMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
