workspace(name = "prysm")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

http_archive(
    name = "bazel_skylib",
    sha256 = "f7be3474d42aae265405a592bb7da8e171919d74c16f082a5457840f06054728",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-skylib/releases/download/1.2.1/bazel-skylib-1.2.1.tar.gz",
        "https://github.com/bazelbuild/bazel-skylib/releases/download/1.2.1/bazel-skylib-1.2.1.tar.gz",
    ],
)

load("@bazel_skylib//:workspace.bzl", "bazel_skylib_workspace")

bazel_skylib_workspace()

http_archive(
    name = "bazel_toolchains",
    sha256 = "8e0633dfb59f704594f19ae996a35650747adc621ada5e8b9fb588f808c89cb0",
    strip_prefix = "bazel-toolchains-3.7.0",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-toolchains/releases/download/3.7.0/bazel-toolchains-3.7.0.tar.gz",
        "https://github.com/bazelbuild/bazel-toolchains/releases/download/3.7.0/bazel-toolchains-3.7.0.tar.gz",
    ],
)

BAZEL_TOOLCHAIN_REF = "056aeaa01900f5050a9fed9b11e2d365a684831a"

BAZEL_TOOLCHAIN_SHA = "93aa940bcaa2bfdd8153d4d029bad1ccc6c0601e29ffff3a23e1d89aba5f61fa"

http_archive(
    name = "com_grail_bazel_toolchain",
    canonical_id = BAZEL_TOOLCHAIN_REF,
    sha256 = BAZEL_TOOLCHAIN_SHA,
    strip_prefix = "bazel-toolchain-{ref}".format(ref = BAZEL_TOOLCHAIN_REF),
    url = "https://github.com/grailbio/bazel-toolchain/archive/{ref}.tar.gz".format(ref = BAZEL_TOOLCHAIN_REF),
)

load("@com_grail_bazel_toolchain//toolchain:deps.bzl", "bazel_toolchain_dependencies")

bazel_toolchain_dependencies()

load("@com_grail_bazel_toolchain//toolchain:rules.bzl", "llvm_toolchain")

llvm_toolchain(
    name = "llvm_11_toolchain",
    cxx_standard = {
        "": "c++14",
    },
    llvm_version = "11.1.0",
    target_settings = {
        "": ["@//:clang11_requested"],
    },
)

load("@llvm_11_toolchain//:toolchains.bzl", llvm_11_register_toolchains = "llvm_register_toolchains")

llvm_11_register_toolchains()

llvm_toolchain(
    name = "llvm_14_toolchain",
    cxx_standard = {
        "": "c++14",
    },
    llvm_version = "14.0.0",
    target_settings = {
        "": ["@//:clang14_requested"],
    },
)

load("@llvm_14_toolchain//:toolchains.bzl", llvm_14_register_toolchains = "llvm_register_toolchains")

llvm_14_register_toolchains()

http_archive(
    name = "aspect_gcc_toolchain",
    patch_args = ["-p1"],
    patches = [
        "@//:third_party/aspect_gcc_toolchain/0001-Expose-target_settings-and-set-std-c-20.patch",
    ],
    sha256 = "e2e12202dd83f67d71101b24554044de25e1625d16b4b56bc453ecaa8f7c6bd0",
    strip_prefix = "aspect-build-gcc-toolchain-ac745d4",
    type = "tar.gz",
    urls = [
        "https://github.com/aspect-build/gcc-toolchain/tarball/ac745d4685e2095cc4f057862800f3f0a473c201",
    ],
)

load("@aspect_gcc_toolchain//toolchain:defs.bzl", "ARCHS", "gcc_register_toolchain")

gcc_register_toolchain(
    name = "gcc10_toolchain_x86_64",
    gcc_version = "10.3.0",
    sysroot_variant = "x86_64",
    target_arch = ARCHS.x86_64,
    target_settings = ["@//:gcc10_requested"],
)

# GCC 11.3.0
gcc_register_toolchain(
    name = "gcc11_x86-64",
    sha256 = "861c1e8ad0a66e4c28e7a1f8319d68080ab0ff8d16a765e65540f1957203a190",
    strip_prefix = "x86-64--glibc--stable-2022.08-1",
    target_arch = "x86_64",
    target_settings = ["@//:gcc11_requested"],
    url = "https://toolchains.bootlin.com/downloads/releases/toolchains/x86-64/tarballs/x86-64--glibc--stable-2022.08-1.tar.bz2",
)

# GCC 12.2.0
gcc_register_toolchain(
    name = "gcc12_x86-64",
    sha256 = "9a09ac03001ef2a6cab391cc658fc2a32730b6a8f25614e97a91b9a79537fe33",
    strip_prefix = "x86-64--glibc--bleeding-edge-2022.08-1",
    target_arch = "x86_64",
    target_settings = ["@//:gcc12_requested"],
    url = "https://toolchains.bootlin.com/downloads/releases/toolchains/x86-64/tarballs/x86-64--glibc--bleeding-edge-2022.08-1.tar.bz2",
)

# gcc_register_toolchain(
#     name = "gcc12_toolchain_x86_64",
#     gcc_version = "12.2.1",
#     sysroot_variant = "x86_64",
#     target_arch = ARCHS.x86_64,
#     target_settings = ["@//:gcc12_requested"],
# )

load("@prysm//tools/cross-toolchain:prysm_toolchains.bzl", "configure_prysm_toolchains")

configure_prysm_toolchains()

load("@prysm//tools/cross-toolchain:rbe_toolchains_config.bzl", "rbe_toolchains_config")

rbe_toolchains_config()

http_archive(
    name = "bazel_gazelle",
    sha256 = "5982e5463f171da99e3bdaeff8c0f48283a7a5f396ec5282910b9e8a49c0dd7e",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.25.0/bazel-gazelle-v0.25.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.25.0/bazel-gazelle-v0.25.0.tar.gz",
    ],
)

http_archive(
    name = "com_github_atlassian_bazel_tools",
    sha256 = "60821f298a7399450b51b9020394904bbad477c18718d2ad6c789f231e5b8b45",
    strip_prefix = "bazel-tools-a2138311856f55add11cd7009a5abc8d4fd6f163",
    urls = ["https://github.com/atlassian/bazel-tools/archive/a2138311856f55add11cd7009a5abc8d4fd6f163.tar.gz"],
)

http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "1f4e59843b61981a96835dc4ac377ad4da9f8c334ebe5e0bb3f58f80c09735f4",
    strip_prefix = "rules_docker-0.19.0",
    urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.19.0/rules_docker-v0.19.0.tar.gz"],
)

http_archive(
    name = "io_bazel_rules_go",
    patch_args = ["-p1"],
    patches = [
        # Expose internals of go_test for custom build transitions.
        "//third_party:io_bazel_rules_go_test.patch",
    ],
    sha256 = "ae013bf35bd23234d1dea46b079f1e05ba74ac0321423830119d3e787ec73483",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.36.0/rules_go-v0.36.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.36.0/rules_go-v0.36.0.zip",
    ],
)

# Override default import in rules_go with special patch until
# https://github.com/gogo/protobuf/pull/582 is merged.
git_repository(
    name = "com_github_gogo_protobuf",
    commit = "b03c65ea87cdc3521ede29f62fe3ce239267c1bc",
    patch_args = ["-p1"],
    patches = [
        "@io_bazel_rules_go//third_party:com_github_gogo_protobuf-gazelle.patch",
        "//third_party:com_github_gogo_protobuf-equal.patch",
    ],
    remote = "https://github.com/gogo/protobuf",
    shallow_since = "1610265707 +0000",
    # gazelle args: -go_prefix github.com/gogo/protobuf -proto legacy
)

load(
    "@io_bazel_rules_docker//repositories:repositories.bzl",
    container_repositories = "repositories",
)

container_repositories()

load(
    "@io_bazel_rules_docker//container:container.bzl",
    "container_pull",
)

container_pull(
    name = "cc_image_base",
    digest = "sha256:41036fc7ed8df0f6addc18484cef0c94a85867508967789f947e11ffd5ff0cc8",
    registry = "gcr.io",
    repository = "distroless/cc",
)

container_pull(
    name = "cc_debug_image_base",
    digest = "sha256:6865ad48467c89c3c3524d4c426f52ad12d9ab7dec31fad31fae69da40eb6445",
    registry = "gcr.io",
    repository = "distroless/cc",
)

container_pull(
    name = "go_image_base",
    digest = "sha256:b9b124f955961599e72630654107a0cf04e08e6fa777fa250b8f840728abd770",
    registry = "gcr.io",
    repository = "distroless/base",
)

container_pull(
    name = "go_debug_image_base",
    digest = "sha256:65668d2b78d25df3d8ccf5a778d017fcaba513b52078c700083eaeef212b9979",
    registry = "gcr.io",
    repository = "distroless/base",
)

container_pull(
    name = "alpine_cc_linux_amd64",
    digest = "sha256:752aa0c9a88461ffc50c5267bb7497ef03a303e38b2c8f7f2ded9bebe5f1f00e",
    registry = "index.docker.io",
    repository = "pinglamb/alpine-glibc",
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains(
    go_version = "1.19.4",
    nogo = "@//:nogo",
)

http_archive(
    name = "io_kubernetes_build",
    sha256 = "b84fbd1173acee9d02a7d3698ad269fdf4f7aa081e9cecd40e012ad0ad8cfa2a",
    strip_prefix = "repo-infra-6537f2101fb432b679f3d103ee729dd8ac5d30a0",
    url = "https://github.com/kubernetes/repo-infra/archive/6537f2101fb432b679f3d103ee729dd8ac5d30a0.tar.gz",
)

http_archive(
    name = "eip3076_spec_tests",
    build_file_content = """
filegroup(
    name = "test_data",
    srcs = glob([
        "**/*.json",
    ]),
    visibility = ["//visibility:public"],
)
    """,
    sha256 = "91434d5fd5e1c6eb7b0174fed2afe25e09bddf00e1e4c431db931b2cee4e7773",
    url = "https://github.com/eth-clients/slashing-protection-interchange-tests/archive/b8413ca42dc92308019d0d4db52c87e9e125c4e9.tar.gz",
)

consensus_spec_version = "v1.3.0-rc.2"

consensus_spec_test_version = "v1.3.0-rc.2-hotfix"

bls_test_version = "v0.1.1"

http_archive(
    name = "consensus_spec_tests_general",
    build_file_content = """
filegroup(
    name = "test_data",
    srcs = glob([
        "**/*.ssz_snappy",
        "**/*.yaml",
    ]),
    visibility = ["//visibility:public"],
)
    """,
    sha256 = "11fed0121d7a79b6da3671ecbd8aa08bcd3fc2316450eeae2d0e065371910ef9",
    url = "https://github.com/ethereum/consensus-spec-tests/releases/download/%s/general.tar.gz" % consensus_spec_test_version,
)

http_archive(
    name = "consensus_spec_tests_minimal",
    build_file_content = """
filegroup(
    name = "test_data",
    srcs = glob([
        "**/*.ssz_snappy",
        "**/*.yaml",
    ]),
    visibility = ["//visibility:public"],
)
    """,
    sha256 = "8fbec5f1fa30fb1194b42230bf75f05690a72cf3c56af5b5ac8aafb94c60fc19",
    url = "https://github.com/ethereum/consensus-spec-tests/releases/download/%s/minimal.tar.gz" % consensus_spec_test_version,
)

http_archive(
    name = "consensus_spec_tests_mainnet",
    build_file_content = """
filegroup(
    name = "test_data",
    srcs = glob([
        "**/*.ssz_snappy",
        "**/*.yaml",
    ]),
    visibility = ["//visibility:public"],
)
    """,
    sha256 = "cb22875dbb67b12f577d12a5aa8074de4c8bd31d74fff45cc110be515c1df847",
    url = "https://github.com/ethereum/consensus-spec-tests/releases/download/%s/mainnet.tar.gz" % consensus_spec_test_version,
)

http_archive(
    name = "consensus_spec",
    build_file_content = """
filegroup(
    name = "spec_data",
    srcs = glob([
        "**/*.yaml",
    ]),
    visibility = ["//visibility:public"],
)
    """,
    sha256 = "c728947ccf8b628691c76cc92f2f59c02f1c7882a463f26365d32c098bb946b8",
    strip_prefix = "consensus-specs-" + consensus_spec_version[1:],
    url = "https://github.com/ethereum/consensus-specs/archive/refs/tags/%s.tar.gz" % consensus_spec_version,
)

http_archive(
    name = "bls_spec_tests",
    build_file_content = """
filegroup(
    name = "test_data",
    srcs = glob([
        "**/*.yaml",
    ]),
    visibility = ["//visibility:public"],
)
    """,
    sha256 = "93c7d006e7c5b882cbd11dc9ec6c5d0e07f4a8c6b27a32f964eb17cf2db9763a",
    url = "https://github.com/ethereum/bls12-381-tests/releases/download/%s/bls_tests_yaml.tar.gz" % bls_test_version,
)

http_archive(
    name = "eth2_networks",
    build_file_content = """
filegroup(
    name = "configs",
    srcs = glob([
        "shared/**/config.yaml",
    ]),
    visibility = ["//visibility:public"],
)
    """,
    sha256 = "82b01a48b143fe0f2fb7fb5f5dd385c1f934335a12d7954f08b1d45d77427b5e",
    strip_prefix = "eth2-networks-674f7a1d01d9c18345456eab76e3871b3df2126b",
    url = "https://github.com/eth-clients/eth2-networks/archive/674f7a1d01d9c18345456eab76e3871b3df2126b.tar.gz",
)

http_archive(
    name = "com_github_bazelbuild_buildtools",
    sha256 = "7a182df18df1debabd9e36ae07c8edfa1378b8424a04561b674d933b965372b3",
    strip_prefix = "buildtools-f2aed9ee205d62d45c55cfabbfd26342f8526862",
    url = "https://github.com/bazelbuild/buildtools/archive/f2aed9ee205d62d45c55cfabbfd26342f8526862.zip",
)

git_repository(
    name = "com_google_protobuf",
    commit = "436bd7880e458532901c58f4d9d1ea23fa7edd52",
    remote = "https://github.com/protocolbuffers/protobuf",
    shallow_since = "1617835118 -0700",
)

# Group the sources of the library so that CMake rule have access to it
all_content = """filegroup(name = "all", srcs = glob(["**"]), visibility = ["//visibility:public"])"""

# External dependencies

http_archive(
    name = "prysm_web_ui",
    build_file_content = """
filegroup(
    name = "site",
    srcs = glob(["**/*"]),
    visibility = ["//visibility:public"],
)
""",
    sha256 = "b2226874526805d64c29e5053fa28e511b57c0860585d6d59777ee81ff4859ca",
    urls = [
        "https://github.com/prysmaticlabs/prysm-web-ui/releases/download/v2.0.2/prysm-web-ui.tar.gz",
    ],
)

load("//:deps.bzl", "prysm_deps")

# gazelle:repository_macro deps.bzl%prysm_deps
prysm_deps()

load("@prysm//third_party/herumi:herumi.bzl", "bls_dependencies")

bls_dependencies()

load("@prysm//testing/endtoend:deps.bzl", "e2e_deps")

e2e_deps()

load(
    "@io_bazel_rules_docker//go:image.bzl",
    _go_image_repos = "repositories",
)

# Golang images
# This is using gcr.io/distroless/base
_go_image_repos()

# CC images
# This is using gcr.io/distroless/base
load(
    "@io_bazel_rules_docker//cc:image.bzl",
    _cc_image_repos = "repositories",
)

_cc_image_repos()

load("@io_bazel_rules_go//extras:embed_data_deps.bzl", "go_embed_data_dependencies")

go_embed_data_dependencies()

load("@com_github_atlassian_bazel_tools//gometalinter:deps.bzl", "gometalinter_dependencies")

gometalinter_dependencies()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

load("@com_github_bazelbuild_buildtools//buildifier:deps.bzl", "buildifier_dependencies")

buildifier_dependencies()

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()

# Do NOT add new go dependencies here! Refer to DEPENDENCIES.md!
