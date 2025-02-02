// Copyright 2021 PingCAP, Inc. Licensed under Apache-2.0.

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MyStruct struct {
	FirstField  string `matched:"first tag" value:"whatever"`
	SecondField string `matched:"second tag" value:"another whatever"`
}

func TestGetFieldTags(t *testing.T) {
	rst := GetFieldsAndTags(MyStruct{}, []string{"matched", "value"})

	assert.Equal(t, rst, []Field{
		{
			Name: "FirstField",
			Tags: map[string]string{
				"matched": "first tag",
				"value":   "whatever",
			},
		},
		{

			Name: "SecondField",
			Tags: map[string]string{
				"matched": "second tag",
				"value":   "another whatever",
			},
		},
	})
}

// // TODO: support nested struct
// func TestGetFieldTags_with_nested_struct(t *testing.T) {}
