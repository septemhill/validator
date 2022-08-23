package validator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleGroupValidator(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy-path",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewSimpleGroupValidator()
			err := v.Validate(tt.args.ctx, nil)
			assert.Equal(t, tt.wantErr, (err != nil))
		})
	}
}
