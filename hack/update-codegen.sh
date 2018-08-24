#!/bin/bash

# Copyright 2018 The Jetstack tarmak contributors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# The only argument this script should ever be called with is '--verify-only'

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE}")/..
CODEGEN_PKG=${CODEGEN_PKG:-$(CDPATH='' cd "${SCRIPT_ROOT}"; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}

"${CODEGEN_PKG}/generate-internal-groups.sh" "deepcopy,defaulter" \
  github.com/jetstack/tarmak/pkg/tarmak/client \
  github.com/jetstack/tarmak/pkg/apis \
  github.com/jetstack/tarmak/pkg/apis \
  "cluster:v1alpha1 tarmak:v1alpha1" \
  --output-base "${GOPATH}/src/" \
  --go-header-file "${SCRIPT_ROOT}/hack/boilerplate/boilerplate.go.txt"

"${CODEGEN_PKG}/generate-internal-groups.sh" "defaulter,deepcopy,conversion,client,informer,lister" \
  github.com/jetstack/tarmak/pkg/wing/client \
  github.com/jetstack/tarmak/pkg/apis \
  github.com/jetstack/tarmak/pkg/apis \
  wing:v1alpha1 \
  --output-base "${GOPATH}/src/" \
  --go-header-file "${SCRIPT_ROOT}/hack/boilerplate/boilerplate.go.txt"
