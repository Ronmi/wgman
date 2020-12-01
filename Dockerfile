# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.

FROM golang as builder

ADD . /wgman
WORKDIR /wgman
RUN go mod download && go build

FROM ronmi/mingo

COPY --from=builder /wgman/wgman /wgman
ENTRYPOINT ["/wgman"]
