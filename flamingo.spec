Name:    flamingo
Version: v0.1.0
Release: 1
Group:   System Environment/Daemons
Summary: Cloud Instance Contextualization Tool
License: Apache 2.0
URL:     http://github.com/tmrts/flamingo
Source0: https://github.com/tmrts/%{name}/releases/download/%{version}/%{name}-%{version}-linux-amd64.tar.gz

Buildroot: %{_tmppath}/%{name}-%{version}-%(%{__id_u} -n)
Packager:  Tamer Tas <contact@tmrts.com>

Requires(pre): nss
Requires(pre): systemd
Requires(pre): libssh2
Requires(pre): iptables
Requires(pre): shadow-utils
Requires(pre): systemd-units

%description
Flamingo is a lightweight contextualization tool that handles
initialization of cloud instances.

Instances in the cloud needs to be initialized(contextualized).
Contextualization includes User, group, network, systemd, and disk configurations.

Flamingo is written in Go and it's main focus points are:

* Speed
* Lightweight
* Extensibility

%prep
%setup -n %{name}-%{version}.x86_64

%install
install -d -m 755 %{buildroot}/bin

install    -m 755 %{_builddir}/%{name}-%{version}.x86_64/bin/%{name} %{buildroot}/bin/
install    -m 755 %{_builddir}/%{name}-%{version}.x86_64/unit %{buildroot}/

ln -sf %{buildroot}/bin/%{name} /bin/%{name}
systemctl enable --system %{buildroot}/unit/%{name}.service

%files
/bin/*
/unit/*

%changelog
* Mon Aug 03 2015 Flamingo <contact@tmrts.com> 0.1.0
- Initial RPM release.
