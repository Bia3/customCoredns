## test build
create a `default.nix` file in this directory
```nix
{
    buildGoModule ? null,
    lib
}:
buildGoModule rec {
    pname = "coredns";
    version = "0.0.1";

    src = ./.;
    vendorHash = "sha256-v2IYZJw2XITksBnEAlMIFDaHbvffBFZlWhaYwl8PCMo=";
    
    meta = with lib; {
        description = "Coredns with plugins for nixos configurations";
        homepage = "https://github.com/bia3/customCoredns";
        license = licenses.mit;
        maintainers = with maintainers; [ bia3 ];
    };
}
```
and ru the following
```shell
nix-build -E 'with import <nixpkgs> { }; callPackage ./default.nix { }'
```