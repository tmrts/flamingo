Name:    flamingo
Version: v0.2.0
Release: 2
Group:   System Environment/Daemons
Summary: Cloud Instance Contextualization Tool
License: Apache 2.0
URL:     http://github.com/tmrts/flamingo
Source0: https://github.com/tmrts/%{name}/archive/%{version}.tar.gz#/%{name}-%{version}.tar.gz

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
%setup -qn %{name}-%{version}

%install
install    -d 755 %{buildroot}/bin
install    -d 755 %{buildroot}/unit

install    -m 755 %{_builddir}/%{name}-%{version}/bin/%{name} %{buildroot}/bin
install    -m 755 %{_builddir}/%{name}-%{version}/unit/%{name}.service %{buildroot}/unit

%post
ln -sf %{buildroot}/bin/%{name} /bin/%{name}
%systemd_post %{name}.service

%files
/bin/*

%changelog
* Wed Aug 12 2015 Flamingo <contact@tmrts.com> 0.2.0
- Openstack, EC2 Cloud Provider Support. Couple of Bug Fixes
* Mon Aug 03 2015 Flamingo <contact@tmrts.com> 0.1.0
- Initial RPM release.
