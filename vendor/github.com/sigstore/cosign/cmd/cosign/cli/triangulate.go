//
// Copyright 2021 The Sigstore Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/sigstore/cosign/pkg/cosign"
)

func Triangulate() *ffcli.Command {
	var (
		flagset = flag.NewFlagSet("cosign triangulate", flag.ExitOnError)
	)
	return &ffcli.Command{
		Name:       "triangulate",
		ShortUsage: "cosign triangulate <image uri>",
		ShortHelp:  "Outputs the located cosign image reference. This is the location cosign stores signatures.",
		FlagSet:    flagset,
		Exec: func(ctx context.Context, args []string) error {
			if len(args) != 1 {
				return flag.ErrHelp
			}
			return MungeCmd(ctx, args[0])
		},
	}
}

func MungeCmd(ctx context.Context, imageRef string) error {
	ref, err := name.ParseReference(imageRef)
	if err != nil {
		return err
	}

	h, err := Digest(ctx, ref)
	if err != nil {
		return err
	}

	sigRepo, err := TargetRepositoryForImage(ref)
	if err != nil {
		return err
	}
	dstRef := cosign.AttachedImageTag(sigRepo, h, cosign.SignatureTagSuffix)

	fmt.Println(dstRef.Name())
	return nil
}