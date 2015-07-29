Name:    flamingo
Version: 0.1.0
Release: 1
Summary: Cloud Instance Contextualization Tool
License: Apache 2.0
URL:     http://github.com/tmrts/flamingo
Group:   System Environment/Daemons
Source0: https://github.com/tmrts/%{name}/releases/download/v%{version}/%%{name}-v%{version}-linux-amd64.tar.gz

Buildroot: %{_tmppath}/%{name}-%{version}-%{release}-%(%{__id_u} -n)
Packager:  Tamer Tas <contact@tmrts.com>

Requires(pre): shadow-utils
Requires(pre): systemd-units
Requires(pre): iptables

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
%setup -n %{name}-v%{version}-linux-amd64

%build
rm -rf %{buildroot}

echo  %{buildroot}

%install
install -d -m 755 %{buildroot}/%{_sbindir}
install    -m 755 %{_builddir}/%{name}-v%{version}-linux-amd64/flamingo %{buildroot}/%{_sbindir}

%clean
rm -rf %{buildroot}

%post
chkconfig --add %{name}

%files
%defattr(-,root,root)
%{_sbindir}/flamingo
%{_defaultdocdir}/doc/*.md

%changelog
    - Initial Spec.
