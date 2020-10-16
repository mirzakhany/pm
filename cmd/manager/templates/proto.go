package templates

const ProtoTmpl = `
syntax = "proto3";

package {{ .Pkg.NamePlural | lower }}V1;

option go_package = "protobuf/{{ .Pkg.NamePlural | lower }};{{ .Pkg.NamePlural | lower }}";

import "google/api/annotations.proto";
import "google/protobuf/empty.proto";
import "google/protobuf/timestamp.proto";

message List{{ .Pkg.NamePlural }}Request {
    int64 limit = 1;
    int64 offset = 2;
}

message List{{ .Pkg.NamePlural }}Response {
    repeated {{ .Pkg.Name }} {{ .Pkg.NamePlural | lower }} = 1;
    int64 total_count = 2;
    int64 limit = 3;
    int64 offset = 4;
}

message Get{{ .Pkg.Name }}Request {
    string uuid = 1;
}

message Create{{ .Pkg.Name }}Request {
    string title = 1;
}

message Update{{ .Pkg.Name }}Request {
    string uuid = 1;
    string title = 2;
}

message Delete{{ .Pkg.Name }}Request {
    string uuid = 1;
}

message {{ .Pkg.Name }} {
    uint64 id = 1;
    string uuid = 2;
    string title = 3;
    google.protobuf.Timestamp created_at = 4;
    google.protobuf.Timestamp updated_at = 5;
}

service {{ .Pkg.Name }}Service {

    // List {{ .Pkg.NamePlural }}
    rpc List{{ .Pkg.NamePlural }} (List{{ .Pkg.NamePlural }}Request) returns (List{{ .Pkg.NamePlural }}Response) {
        option (google.api.http) = {
            get: "/v1/{{ .Pkg.NamePlural | lower }}"
        };
    }
    // Get {{ .Pkg.Name }}
    rpc Get{{ .Pkg.Name }} (Get{{ .Pkg.Name }}Request) returns ({{ .Pkg.Name }}) {
        option (google.api.http) = {
          get: "/v1/{{ .Pkg.NamePlural | lower }}/{uuid}"
        };
    }

    // Create {{ .Pkg.Name }} object request
    rpc Create{{ .Pkg.Name }} (Create{{ .Pkg.Name }}Request) returns ({{ .Pkg.Name }}) {
        option (google.api.http) = {
            post: "/v1/{{ .Pkg.NamePlural | lower }}"
            body: "*"
        };
    }

    // Update {{ .Pkg.Name }} object request
    rpc Update{{ .Pkg.Name }} (Update{{ .Pkg.Name }}Request) returns ({{ .Pkg.Name }}) {
        option (google.api.http) = {
            put: "/v1/{{ .Pkg.NamePlural | lower }}/{uuid}"
            body: "*"
        };
    }

    // Delete {{ .Pkg.Name }} object request
    rpc Delete{{ .Pkg.Name }} (Delete{{ .Pkg.Name }}Request) returns (google.protobuf.Empty) {
        option (google.api.http) = {
          delete: "/v1/{{ .Pkg.NamePlural | lower }}/{uuid}"
        };
    }
}
`
