load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")  # gazelle:keep

lighthouse_version = "v2.4.0"
lighthouse_archive_name = "lighthouse-%s-x86_64-unknown-linux-gnu-portable.tar.gz" % lighthouse_version

def e2e_deps():
    http_archive(
        name = "web3signer",
        urls = ["https://artifacts.consensys.net/public/web3signer/raw/names/web3signer.tar.gz/versions/22.7.0/web3signer-22.7.0.tar.gz"],
        sha256 = "a4b7c6261776c651bc9016c73d28e99190d97a77c2d661715ae1902eedffb1c0",
        build_file = "@prysm//testing/endtoend:web3signer.BUILD",
        strip_prefix = "web3signer-22.7.0",
    )

    http_archive(
        name = "lighthouse",
        sha256 = "bfa32f62281f57d234881278ed9ca4712a946907b9e6ad6cd7f6e3a3bdea5364",
        build_file = "@prysm//testing/endtoend:lighthouse.BUILD",
        url = ("https://github.com/sigp/lighthouse/releases/download/%s/" + lighthouse_archive_name) % lighthouse_version,
    )
