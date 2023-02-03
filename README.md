# tar256sum
Repeatable tarball checksums in 100 lines of go.

- Doesn't write to disk
- Fully streaming - stable memory footprint
- Resilient from zip-bomb attacks (see `--max-decompress`)
- Stable output

## Why

GitHub recently [broke](https://github.com/orgs/community/discussions/45830) many build pipelines [by accident](https://github.blog/changelog/2023-01-30-git-archive-checksums-may-change/). This project explores an alternative checksum method that might be more resilient.

## Install

Nix:

```
nix profile install github:andrewhamon/tar256sum
```

Go get:

```
go install github.com/andrewhamon/tar256sum
```

## Usage

```
cat archive.tar.gz | tar256sum
```

## Demo

1. Check out a sizeable git repo:
    ```
    git clone https://github.com/NixOS/nixpkgs.git
    cd nixpkgs
    ```

2. Have an older version of git handy:
    ```
    oldgit --version
    # git version 2.36.2

    git --version
    # git version 2.38.3
    ```

3. Compare `sha256sum` for two archives of the same commit:
    ```

    git archive --format tar.gz dbae6eb51edb8afe281e995eff341be07fc43247 | sha256sum
    # 8d88969fcaf813e4d4c2f1d14f26ad45a2c35108d5419a31001b04c34cad3579  -

    oldgit archive --format tar.gz dbae6eb51edb8afe281e995eff341be07fc43247 | sha256sum
    # f1f69372dbb92c00a16e7f73b03d26d7d0462864df7ab854061952be7976e02c  -
    ```

    Observe the non-repeatability. Feel sad.

4. Try again with tar256sum:
    ```
    git archive --format tar.gz dbae6eb51edb8afe281e995eff341be07fc43247 | tar256sum
    # 1d4d42cf4f450f7dc3c4d071d5ce684029a09b50844d472104405bfa1bfd3efc  -


    oldgit archive --format tar.gz dbae6eb51edb8afe281e995eff341be07fc43247 | tar256sum
    # 1d4d42cf4f450f7dc3c4d071d5ce684029a09b50844d472104405bfa1bfd3efc  -
    ```

    Hooray, a stable result!

## How does this work?

- for each tar entry:
  - checksum the entry header
  - checksum the contents of the entry
  - store this pair of checksums (in memory)
- sort these pairs and checksum them to produce final result

## Why not `cat archive.tar.gz | gunzip | sha256sum`

You know, I'm starting to ask myself the same question. Obviously there is some
zip bomb risk piping to gunzip (is there a flag for that with gunzip?). But I can
only produce different git archive results with compression -- plain tar seems
more stable across versions, so this program could perhaps be made way simpler
by simply checking the raw tar without sorting the entries.

**If you know of two equivalent git tar archives (`git archive --format tar`)
that have different checksums, I would love to know about it.**
