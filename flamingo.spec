Name:    flamingo
Version: 0.2.0
Release: 3
Group:   System Environment/Daemons
Summary: Cloud Instance Contextualization Tool
License: Apache 2.0
URL:     http://github.com/tmrts/flamingo
Source0: https://github.com/tmrts/%{name}/archive/v%{version}.tar.gz#/%{name}-v%{version}.tar.gz

Buildroot: %{_tmppath}/%{name}-%{version}-%(%{__id_u} -n)
Packager:  Tamer Tas <contact@tmrts.com>

Requires(post): nss
Requires(post): libssh2
Requires(post): iptables
Requires(post): shadow-utils
Requires(post): systemd
BuildRequires: systemd

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
install -d -m 755 %{buildroot}/bin
install -d -m 755 %{buildroot}/etc/systemd/system/

install -m 700 %{_builddir}/%{name}-%{version}/bin/%{name} %{buildroot}/bin
install -m 700 %{_builddir}/%{name}-%{version}/unit/%{name}.service %{buildroot}/etc/systemd/system/

%clean
rm -rf %{buildroot}

%post
systemctl enable --system /etc/systemd/system/%{name}.service

%files
/bin/*
/etc/systemd/system/*

%changelog
* Wed Aug 12 2015 Flamingo <contact@tmrts.com> 0.2.0
- Openstack, EC2 Cloud Provider Support. Couple of Bug Fixes
* Mon Aug 03 2015 Flamingo <contact@tmrts.com> 0.1.0
- Initial RPM release.
