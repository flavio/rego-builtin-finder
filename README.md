Rego policies can be built into [WebAssembly modules](https://www.openpolicyagent.org/docs/latest/wasm/).

The Rego language offers a series of [built-in functions](https://www.openpolicyagent.org/docs/latest/policy-reference/#built-in-functions)
that can be used to write policies.

During the process of compiling a Rego file to a WebAssembly module, some of
the built-in functions used by the policy are automatically translated to
WebAssembly. Others instead have to be provided at runtime by the host
evaluating the WebAssembly; these are called "SDK dependent" built-ins.

This cli tool scans Rego files and reports all the SDK dependent built-ins
that are used.

## Usage

The tool can scan a directory recursively, analyzing all the Rego files it
finds.

**Note well:** Rego test files are automatically ignored.

A quick example:

```console
$ git clone git@github.com:open-policy-agent/library.git opa-library
$ rego-builtin-finder ./opa-library
Rego files analyzed: 33
List of builtins that have to be provided by the SDK

NAME          	OCCURRENCES
sprintf       	5
yaml.unmarshal	1
http.send     	1
cast_array    	1
```

The tool can also be pointed to a specific Rego file to be analyzed:

```console
$ rego-builtin-finder ./opa-library/kubernetes/mutating-admission/test_mutation.rego
Rego files analyzed: 1
List of builtins that have to be provided by the SDK

NAME   	OCCURRENCES
trace  	1
sprintf	1
```

The tool can print additional debug information via the `-debug` flag.
