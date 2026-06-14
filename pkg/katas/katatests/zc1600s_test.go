// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1600(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — bwrap for sandboxed exec",
			input:    `bwrap --ro-bind / / --unshare-user --uid 1000 /bin/sh`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — no chroot call",
			input:    `mount --bind /src /dst`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — chroot /var/sandbox /bin/sh",
			input: `chroot /var/sandbox /bin/sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1600",
					Message: "`chroot` without `--userspec=` runs the inner command as uid 0. Pass `--userspec=USER:GROUP` to drop privileges, or use `bwrap` / `firejail` for user-namespace isolation.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — chroot $ROOT sh",
			input: `chroot $ROOT sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1600",
					Message: "`chroot` without `--userspec=` runs the inner command as uid 0. Pass `--userspec=USER:GROUP` to drop privileges, or use `bwrap` / `firejail` for user-namespace isolation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1600")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1601(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ethtool wol d (disable)",
			input:    `ethtool -s eth0 wol d`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ethtool setting different knob",
			input:    `ethtool -s eth0 autoneg on`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ethtool -s eth0 wol g",
			input: `ethtool -s eth0 wol g`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1601",
					Message: "`ethtool -s eth0 wol g` enables Wake-on-LAN — the NIC powers the host on before firewall rules load. Keep `wol d` unless a documented operational need requires g.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ethtool -s $IF wol ubg",
			input: `ethtool -s $IF wol ubg`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1601",
					Message: "`ethtool -s $IF wol ubg` enables Wake-on-LAN — the NIC powers the host on before firewall rules load. Keep `wol d` unless a documented operational need requires ubg.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1601")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1602(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — setopt NO_NOMATCH",
			input:    `setopt NO_NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — setopt EXTENDED_GLOB",
			input:    `setopt EXTENDED_GLOB`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — setopt KSH_ARRAYS",
			input: `setopt KSH_ARRAYS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1602",
					Message: "`setopt KSH_ARRAYS` flips Zsh core semantics for the whole shell — pre-existing code silently misbehaves. Scope with `emulate -L ksh` / `emulate -L sh` inside the function that needs the mode.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — setopt shwordsplit (lowercase, no underscore)",
			input: `setopt shwordsplit`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1602",
					Message: "`setopt shwordsplit` flips Zsh core semantics for the whole shell — pre-existing code silently misbehaves. Scope with `emulate -L ksh` / `emulate -L sh` inside the function that needs the mode.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1602")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1603(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — gdb on core file",
			input:    `gdb /usr/bin/app /var/lib/cores/app.core`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — coredumpctl",
			input:    `coredumpctl debug myapp`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — gdb -p 1234",
			input: `gdb -p 1234`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1603",
					Message: "`gdb -p PID` attaches via ptrace — memory, registers, env, and stack of the target are readable. Use `coredumpctl` on a captured core, not a live attach from a script.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ltrace -p $PID",
			input: `ltrace -p $PID`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1603",
					Message: "`ltrace -p PID` attaches via ptrace — memory, registers, env, and stack of the target are readable. Use `coredumpctl` on a captured core, not a live attach from a script.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1603")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1604(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — source explicit file",
			input:    `source /etc/bashrc`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — source variable path (no glob)",
			input:    `source $HOME/dotfiles/common.sh`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — source /etc/profile.d/*.sh",
			input: `source /etc/profile.d/*.sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1604",
					Message: "`source /etc/profile.d/*.sh` loads every matched file. One attacker-writable match is arbitrary code execution. Use explicit filenames.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — . $HOME/dotfiles/*.sh",
			input: `. $HOME/dotfiles/*.sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1604",
					Message: "`. $HOME/dotfiles/*.sh` loads every matched file. One attacker-writable match is arbitrary code execution. Use explicit filenames.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			// A process substitution is dynamic command output, not a
			// filename glob; a `*` inside its body is not a match pattern.
			name:     "valid — source process substitution with glob in body",
			input:    `source <(zstyle -L ':z:*' | awk '{print $2}')`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — source command substitution",
			input:    `source $(find . -name '*.sh')`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1604")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1605(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — read-only debugfs",
			input:    `debugfs $DEV`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — debugfs -R command (read-only)",
			input:    `debugfs -R "stat foo.txt" $DEV`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — debugfs -w $DEV",
			input: `debugfs -w $DEV`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1605",
					Message: "`debugfs -w` writes to the filesystem outside the kernel's normal path — journal bypassed, locks ignored. Keep it as an interactive rescue tool, not a script path.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — debugfs -w /dev/loop0",
			input: `debugfs -w /dev/loop0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1605",
					Message: "`debugfs -w` writes to the filesystem outside the kernel's normal path — journal bypassed, locks ignored. Keep it as an interactive rescue tool, not a script path.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1605")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1606(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — mkdir -m 700",
			input:    `mkdir -m 700 /root/data`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — mkdir -m 755",
			input:    `mkdir -m 755 /opt/app`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — mkdir -m 1777 (sticky)",
			input:    `mkdir -m 1777 /tmp/shared`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — mkdir -m 777",
			input: `mkdir -m 777 /tmp/shared`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1606",
					Message: "`mkdir -m 777` creates a world-writable path without the sticky bit — TOCTOU symlink-attack ground. Use `-m 700` / `-m 750`, or `-m 1777` if a shared sticky dir is actually needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — install -m 666",
			input: `install -m 666 file /tmp/x`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1606",
					Message: "`install -m 666` creates a world-writable path without the sticky bit — TOCTOU symlink-attack ground. Use `-m 700` / `-m 750`, or `-m 1777` if a shared sticky dir is actually needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1606")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1607(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — safe.directory scoped to a path",
			input:    `git config --global safe.directory /workspace/repo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — unrelated git config",
			input:    `git config user.email "me@example.com"`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — git config safe.directory '*'",
			input: `git config --global safe.directory '*'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1607",
					Message: "`git config safe.directory '*'` trusts every directory — defeats CVE-2022-24765 protection. List specific paths, or fix the ownership mismatch.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — git -c safe.directory=* status",
			input: `git -c safe.directory=* status`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1607",
					Message: "`git config safe.directory '*'` trusts every directory — defeats CVE-2022-24765 protection. List specific paths, or fix the ownership mismatch.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1607")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1608(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — find -exec rm {}",
			input:    `find . -type f -exec rm {} \;`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — find -exec sh -c 'positional'",
			input:    `find . -exec sh -c 'grep pat "$1"' _ {} \;`,
			expected: []katas.Violation{},
		},
		{
			name:  `invalid — find -exec sh -c 'echo {}'`,
			input: `find . -exec sh -c 'echo {}' \;`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1608",
					Message: "`find -exec sh -c '... {} ...'` interpolates filenames into the shell script — metacharacters break out. Pass `{}` as a positional arg: `sh -c '... \"$1\"' _ {} \\;`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  `invalid — find -exec bash -c "grep X {}"`,
			input: `find . -exec bash -c "grep X {}" \;`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1608",
					Message: "`find -exec sh -c '... {} ...'` interpolates filenames into the shell script — metacharacters break out. Pass `{}` as a positional arg: `sh -c '... \"$1\"' _ {} \\;`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1608")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1609(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — aa-enforce reapplies enforcement",
			input:    `aa-enforce /etc/apparmor.d/usr.bin.firefox`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — apparmor_parser -r (reload, not remove)",
			input:    `apparmor_parser -r /etc/apparmor.d/profile`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — aa-disable",
			input: `aa-disable /etc/apparmor.d/usr.bin.firefox`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1609",
					Message: "`aa-disable` disables or softens the AppArmor profile — the confined process loses MAC restrictions. Review the profile's intent before disabling in automation.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — aa-complain",
			input: `aa-complain /etc/apparmor.d/usr.bin.firefox`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1609",
					Message: "`aa-complain` disables or softens the AppArmor profile — the confined process loses MAC restrictions. Review the profile's intent before disabling in automation.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — apparmor_parser -R",
			input: `apparmor_parser -R /etc/apparmor.d/profile`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1609",
					Message: "`apparmor_parser -R` removes the AppArmor profile from the kernel — the confined process loses MAC restrictions. Review the profile's intent before removing in automation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1609")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1610(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — curl to temp path",
			input:    `curl -fsSL -o /tmp/download.tar.gz https://example.com/x.tar.gz`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — wget to user home",
			input:    `wget -O $HOME/.local/bin/tool https://example.com/tool`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — curl -o /etc/config",
			input: `curl -fsSL -o /etc/myapp/config.yaml https://example.com/config`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1610",
					Message: "`curl -o /etc/myapp/config.yaml` writes an HTTP response straight into a system path — a compromised URL replaces the target. Download to a temp file, verify, then `install` into place.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — wget -O /usr/local/bin/tool",
			input: `wget -O /usr/local/bin/tool https://example.com/tool`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1610",
					Message: "`wget -O /usr/local/bin/tool` writes an HTTP response straight into a system path — a compromised URL replaces the target. Download to a temp file, verify, then `install` into place.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — wget -O /lib/x.so",
			input: `wget -O /lib/evil.so https://example.com/evil.so`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1610",
					Message: "`wget -O /lib/evil.so` writes an HTTP response straight into a system path — a compromised URL replaces the target. Download to a temp file, verify, then `install` into place.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1610")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1611(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — Zsh ${(U)var}",
			input:    `echo "${(U)var}"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — plain expansion",
			input:    `echo "${var}"`,
			expected: []katas.Violation{},
		},
		{
			name:  `invalid — echo "${var^^}"`,
			input: `echo "${var^^}"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1611",
					Message: "`${var^^}` / `${var,,}` — prefer Zsh `${(U)var}` / `${(L)var}` for case conversion.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  `invalid — print -r "${name,,}"`,
			input: `print -r -- "${name,,}"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1611",
					Message: "`${var^^}` / `${var,,}` — prefer Zsh `${(U)var}` / `${(L)var}` for case conversion.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1611")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1612(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — sysctl tightening ptrace_scope",
			input:    `sysctl -w kernel.yama.ptrace_scope=3`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — unrelated sysctl",
			input:    `sysctl -w vm.swappiness=10`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — sysctl -w kernel.yama.ptrace_scope=0",
			input: `sysctl -w kernel.yama.ptrace_scope=0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1612",
					Message: "`sysctl ... kernel.yama.ptrace_scope=0` disables YAMA ptrace scope (lets any process attach) — defense-in-depth loss. Leave the default unless a measured need justifies it.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — sysctl -w kernel.kptr_restrict=0",
			input: `sysctl -w kernel.kptr_restrict=0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1612",
					Message: "`sysctl ... kernel.kptr_restrict=0` disables kernel pointer restriction (leaks kptrs to /proc) — defense-in-depth loss. Leave the default unless a measured need justifies it.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1612")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1613(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ssh-keygen on host key",
			input:    `ssh-keygen -l -f /etc/ssh/ssh_host_rsa_key`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — cat public host key",
			input:    `cat /etc/ssh/ssh_host_rsa_key.pub`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — cat /etc/ssh/ssh_host_ed25519_key",
			input: `cat /etc/ssh/ssh_host_ed25519_key`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1613",
					Message: "Reading `/etc/ssh/ssh_host_ed25519_key` through a text tool copies private-key material into the process and often into logs / scrollback. Use `ssh-keygen -l -f` for metadata, or pass the path directly to the consumer.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — grep PRIVATE $HOME/.ssh/id_rsa",
			input: `grep PRIVATE $HOME/.ssh/id_rsa`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1613",
					Message: "Reading `$HOME/.ssh/id_rsa` through a text tool copies private-key material into the process and often into logs / scrollback. Use `ssh-keygen -l -f` for metadata, or pass the path directly to the consumer.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1613")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1614(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — expect driving a non-auth dialog",
			input:    `expect -c 'spawn lftp host; expect lftp; send "ls\r"; interact'`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — no expect in use",
			input:    `ssh -i key host cmd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — expect with password",
			input: `expect -c 'spawn ssh user@host; expect password; send "s3cret\r"; interact'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1614",
					Message: "`expect` script contains `password` / `passphrase` — the full argv lands in `ps` and audit logs. Switch to key-based auth, or read the credential from a protected file the expect script opens.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — expect with passphrase",
			input: `expect -c 'spawn ssh-keygen -p -f key; expect passphrase; send "x\r"'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1614",
					Message: "`expect` script contains `password` / `passphrase` — the full argv lands in `ps` and audit logs. Switch to key-based auth, or read the credential from a protected file the expect script opens.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1614")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1615(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — plain date",
			input:    `date "+%Y-%m-%d"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — date +%s only (ZC1119 handles)",
			input:    `date "+%s"`,
			expected: []katas.Violation{},
		},
		{
			name:  `invalid — date "+%s.%N"`,
			input: `date "+%s.%N"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1615",
					Message: "`date \"+%s.%N\"` forks for sub-second time. Use Zsh `$EPOCHREALTIME` / `$epochtime` from `zmodload zsh/datetime`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  `invalid — date "+%s%N"`,
			input: `date "+%s%N"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1615",
					Message: "`date \"+%s%N\"` forks for sub-second time. Use Zsh `$EPOCHREALTIME` / `$epochtime` from `zmodload zsh/datetime`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1615")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1616(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — fsfreeze -u (unfreeze)",
			input:    `fsfreeze -u /mnt/backup`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — different command",
			input:    `mount /mnt/backup`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — fsfreeze -f /mnt/backup",
			input: `fsfreeze -f /mnt/backup`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1616",
					Message: "`fsfreeze -f` freezes the mountpoint — every write hangs until `fsfreeze -u` runs. Wrap the call in `trap 'fsfreeze -u PATH' EXIT` so the thaw fires even on failure.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — fsfreeze -f $ROOTFS",
			input: `fsfreeze -f $ROOTFS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1616",
					Message: "`fsfreeze -f` freezes the mountpoint — every write hangs until `fsfreeze -u` runs. Wrap the call in `trap 'fsfreeze -u PATH' EXIT` so the thaw fires even on failure.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1616")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1617(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — xargs -P 4",
			input:    `xargs -P 4 -n 1 echo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — xargs without -P",
			input:    `xargs -n 10 echo`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — xargs -P 0",
			input: `xargs -P 0 -n 1 echo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1617",
					Message: "`xargs -P 0` spawns one child per input line — CPU / FD / memory exhaustion risk. Use `-P $(nproc)` or an explicit cap.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — xargs -P0 (joined)",
			input: `xargs -P0 -n1 echo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1617",
					Message: "`xargs -P 0` spawns one child per input line — CPU / FD / memory exhaustion risk. Use `-P $(nproc)` or an explicit cap.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1617")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1618(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — git commit -m (no skip)",
			input:    `git commit -m "msg"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — git push --dry-run",
			input:    `git push --dry-run origin main`,
			expected: []katas.Violation{},
		},
		{
			name:  `invalid — git commit --no-verify`,
			input: `git commit --no-verify -m "msg"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1618",
					Message: "`git commit --no-verify` skips pre-commit / commit-msg hooks — lint, test, and secret-scan checks do not run. Reserve for emergencies; scripts should pass the hooks.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  `invalid — git push --no-verify`,
			input: `git push --no-verify origin main`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1618",
					Message: "`git push --no-verify` skips pre-push / commit-msg hooks — lint, test, and secret-scan checks do not run. Reserve for emergencies; scripts should pass the hooks.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  `invalid — git commit -n -m`,
			input: `git commit -n -m "msg"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1618",
					Message: "`git commit -n` skips pre-commit / commit-msg hooks — lint, test, and secret-scan checks do not run. Reserve for emergencies; scripts should pass the hooks.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1618")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1619(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — nfs with nosuid,nodev",
			input:    `mount -t nfs -o rw,nosuid,nodev host:/export /mnt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — local ext4",
			input:    `mount -t ext4 /dev/nvme0n1p1 /data`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — nfs without nosuid/nodev",
			input: `mount -t nfs -o rw host:/export /mnt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1619",
					Message: "`mount -t nfs` without nosuid,nodev — a hostile server can plant setuid binaries or device nodes that the client kernel honors. Add `nosuid,nodev` to the `-o` options.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — cifs with only nosuid",
			input: `mount -t cifs -o username=foo,nosuid //host/share /mnt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1619",
					Message: "`mount -t cifs` without nodev — a hostile server can plant setuid binaries or device nodes that the client kernel honors. Add `nosuid,nodev` to the `-o` options.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1619")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1620(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — tee on log file",
			input:    `tee -a /var/log/app.log`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — tee to tmp staging",
			input:    `tee /tmp/sudoers.new`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — tee /etc/sudoers",
			input: `tee /etc/sudoers`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1620",
					Message: "`tee /etc/sudoers` writes without syntax validation — a typo locks everyone out of sudo. Pipe through `visudo -cf /dev/stdin` or stage in a temp file and `visudo -cf` before `mv`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — tee -a /etc/sudoers.d/custom",
			input: `tee -a /etc/sudoers.d/custom`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1620",
					Message: "`tee /etc/sudoers.d/custom` writes without syntax validation — a typo locks everyone out of sudo. Pipe through `visudo -cf /dev/stdin` or stage in a temp file and `visudo -cf` before `mv`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1620")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1621(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — default tmux",
			input:    `tmux new-session -d`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — tmux -S in XDG_RUNTIME_DIR",
			input:    `tmux -S $XDG_RUNTIME_DIR/tmux-main new-session`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — tmux -S /tmp/shared",
			input: `tmux -S /tmp/shared new-session`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1621",
					Message: "`tmux -S /tmp/shared` places the socket in a world-traversable directory — any local user who can read the socket can attach the session. Use `$XDG_RUNTIME_DIR` or a 0700-scoped parent dir.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — tmux -S /var/tmp/pair",
			input: `tmux -S /var/tmp/pair attach`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1621",
					Message: "`tmux -S /var/tmp/pair` places the socket in a world-traversable directory — any local user who can read the socket can attach the session. Use `$XDG_RUNTIME_DIR` or a 0700-scoped parent dir.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1621")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1622(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — Zsh flag form",
			input:    `echo "${(U)var}"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — plain expansion",
			input:    `print -r -- "${var}"`,
			expected: []katas.Violation{},
		},
		{
			name:  `invalid — echo "${var@U}"`,
			input: `echo "${var@U}"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1622",
					Message: "`${var@U}` — prefer Zsh `${(X)var}` parameter-expansion flags (e.g. `${(U)var}` for uppercase).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  `invalid — print "${path@Q}"`,
			input: `print -r -- "${path@Q}"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1622",
					Message: "`${var@Q}` — prefer Zsh `${(X)var}` parameter-expansion flags (e.g. `${(U)var}` for uppercase).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1622")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1623(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — kill -TERM",
			input:    `kill -TERM $PID`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — kill -CONT (resume)",
			input:    `kill -CONT $PID`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — kill -STOP $PID",
			input: `kill -STOP $PID`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1623",
					Message: "`kill -STOP` halts the target until SIGCONT arrives. Pair every STOP with `trap \"kill -CONT PID\" EXIT` so the resume fires even on failure.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — kill -s STOP $PID",
			input: `kill -s STOP $PID`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1623",
					Message: "`kill -STOP` halts the target until SIGCONT arrives. Pair every STOP with `trap \"kill -CONT PID\" EXIT` so the resume fires even on failure.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — pkill -19",
			input: `pkill -19 slowproc`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1623",
					Message: "`kill -STOP` halts the target until SIGCONT arrives. Pair every STOP with `trap \"kill -CONT PID\" EXIT` so the resume fires even on failure.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1623")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1624(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — az login (interactive)",
			input:    `az login`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — az login --service-principal with federated token",
			input:    `az login --service-principal -u appid -t tenantid`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — az login -p $SECRET",
			input: `az login --service-principal -u appid -p $SECRET`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1624",
					Message: "`az login -p` puts the SP password in argv — visible in `ps` / `/proc/<pid>/cmdline`. Use federated-token OIDC, managed identity, or `AZURE_PASSWORD` via a protected env var.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — az login --password literal",
			input: `az login --password hunter2 -u user`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1624",
					Message: "`az login --password` puts the SP password in argv — visible in `ps` / `/proc/<pid>/cmdline`. Use federated-token OIDC, managed identity, or `AZURE_PASSWORD` via a protected env var.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1624")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1625(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — rm -rf scoped path",
			input:    `rm -rf /tmp/staging`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — rm with --preserve-root=all",
			input:    `rm -rf --preserve-root=all $TARGET`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — rm -rf --no-preserve-root /",
			input: `rm -rf --no-preserve-root /`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1625",
					Message: "`rm --no-preserve-root` disables the GNU safeguard against `rm -rf /`. Remove the flag; if a specific path needs deletion, list it explicitly.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — rm -rf --no-preserve-root $TARGET",
			input: `rm -rf --no-preserve-root $TARGET`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1625",
					Message: "`rm --no-preserve-root` disables the GNU safeguard against `rm -rf /`. Remove the flag; if a specific path needs deletion, list it explicitly.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1625")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1626(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — --set non-secret",
			input:    `helm install myapp chart --set replicas=3`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — -f values.yaml",
			input:    `helm install myapp chart -f /secure/values.yaml`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — --set-file points at path",
			input:    `helm install myapp chart --set-file db.password=/run/secrets/db`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — helm install --set password=...",
			input: `helm install myapp chart --set password=s3cret`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1626",
					Message: "`helm install --set password=s3cret` places a secret value in argv — readable via `ps`. Use `-f values.yaml` or `--set-file password=PATH`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — helm upgrade --set-string token=$TOKEN",
			input: `helm upgrade myapp chart --set-string token=$TOKEN`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1626",
					Message: "`helm upgrade --set-string token=$TOKEN` places a secret value in argv — readable via `ps`. Use `-f values.yaml` or `--set-file token=PATH`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1626")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1627(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — crontab from /etc",
			input:    `crontab /etc/cron.install.conf`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — crontab $HOME path",
			input:    `crontab $HOME/.crontab`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — crontab /tmp/newcron",
			input: `crontab /tmp/newcron`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1627",
					Message: "`crontab /tmp/newcron` reads cron rules from a world-traversable path — a concurrent local user can substitute the file between write and read. Stage the file in `$XDG_RUNTIME_DIR/` or `mktemp -d`, or pipe via `crontab -`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — crontab -u bob /tmp/evil",
			input: `crontab -u bob /tmp/evil`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1627",
					Message: "`crontab /tmp/evil` reads cron rules from a world-traversable path — a concurrent local user can substitute the file between write and read. Stage the file in `$XDG_RUNTIME_DIR/` or `mktemp -d`, or pipe via `crontab -`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1627")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1628(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — plain modprobe",
			input:    `modprobe nvme`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — lsmod",
			input:    `lsmod`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — insmod module.ko",
			input: `insmod evilmod.ko`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1628",
					Message: "`insmod` loads a kernel module bypassing depmod / blacklist — prefer `modprobe MODNAME` so system policy and signature checks apply.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — modprobe -f evilmod",
			input: `modprobe -f evilmod`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1628",
					Message: "`modprobe -f` ignores version-magic and kernel-mismatch checks — a mismatched module can crash or compromise the kernel. Drop the flag and fix the underlying version mismatch.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1628")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1629(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — rsync with explicit path",
			input:    `rsync -a --rsync-path=/usr/bin/rsync src host:dst`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — rsync with no path override",
			input:    `rsync -a src host:dst`,
			expected: []katas.Violation{},
		},
		{
			name:  `invalid — rsync --rsync-path="sudo rsync"`,
			input: `rsync -a --rsync-path="sudo rsync" src host:dst`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1629",
					Message: "`rsync --rsync-path='sudo rsync'` runs remote rsync under privilege escalation. Use a scoped sudoers rule on the remote host and keep the path explicit (`/usr/bin/rsync`).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  `invalid — rsync --rsync-path="doas rsync"`,
			input: `rsync -a --rsync-path="doas rsync" src host:dst`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1629",
					Message: "`rsync --rsync-path='doas rsync'` runs remote rsync under privilege escalation. Use a scoped sudoers rule on the remote host and keep the path explicit (`/usr/bin/rsync`).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1629")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1630(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — localhost bind",
			input:    `php -S 127.0.0.1:8000 -t public`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — php script (not -S)",
			input:    `php artisan migrate`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — php -S 0.0.0.0:8000",
			input: `php -S 0.0.0.0:8000`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1630",
					Message: "`php -S 0.0.0.0:8000` binds the dev server to every interface — unauthenticated access to the working directory. Use `127.0.0.1:PORT` locally, nginx / caddy for external exposure.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — php -S [::]:8080",
			input: `php -S [::]:8080 -t public`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1630",
					Message: "`php -S [::]:8080` binds the dev server to every interface — unauthenticated access to the working directory. Use `127.0.0.1:PORT` locally, nginx / caddy for external exposure.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1630")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1631(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — passin env:VAR",
			input:    `openssl pkcs12 -in f.p12 -passin env:PASS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — passin file:path",
			input:    `openssl pkcs12 -in f.p12 -passin file:/run/secrets/p`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — passin stdin",
			input:    `openssl req -passin stdin`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — passin pass:LITERAL",
			input: `openssl pkcs12 -in f.p12 -passin pass:hunter2 -nocerts`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1631",
					Message: "`openssl -passin pass:hunter2` puts the password in argv — visible via `ps`. Use `env:VARNAME`, `file:PATH`, `fd:N`, or `stdin`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — passout pass:X",
			input: `openssl genrsa -passout pass:X 2048`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1631",
					Message: "`openssl -passout pass:X` puts the password in argv — visible via `ps`. Use `env:VARNAME`, `file:PATH`, `fd:N`, or `stdin`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1631")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1632(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — unrelated command",
			input:    `rm /tmp/staging.log`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — shred -u FILE",
			input: `shred -u /tmp/secret.key`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1632",
					Message: "`shred` may not overwrite original blocks on ext4/btrfs/zfs. For guaranteed erasure, use full-disk encryption with key destruction, or `blkdiscard` when retiring an SSD.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — shred -n 3 file",
			input: `shred -n 3 /var/log/secret.log`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1632",
					Message: "`shred` may not overwrite original blocks on ext4/btrfs/zfs. For guaranteed erasure, use full-disk encryption with key destruction, or `blkdiscard` when retiring an SSD.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1632")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1633(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — passphrase-file",
			input:    `gpg -d --passphrase-file /run/secrets/gpg file.gpg`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — passphrase-fd",
			input:    `gpg -d --passphrase-fd 0 file.gpg`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — gpg -d --passphrase 'secret'",
			input: `gpg -d --passphrase 'secret' file.gpg`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1633",
					Message: "`gpg --passphrase` puts the passphrase in argv — visible via `ps`. Use `--passphrase-file`, `--passphrase-fd`, or `--pinentry-mode=loopback` with the value on stdin.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — gpg2 -d --passphrase $PW",
			input: `gpg2 -d --passphrase $PW file.gpg`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1633",
					Message: "`gpg --passphrase` puts the passphrase in argv — visible via `ps`. Use `--passphrase-file`, `--passphrase-fd`, or `--pinentry-mode=loopback` with the value on stdin.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1633")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1634(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — umask 022",
			input:    `umask 022`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — umask 077",
			input:    `umask 077`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — umask 002 (group-write collab)",
			input:    `umask 002`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — umask 111",
			input: `umask 111`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1634",
					Message: "`umask 111` leaves world-write on new files — the \"other\" digit must be `2`/`3`/`6`/`7` to mask the write bit. Use `022` for public, `077` for secrets.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — umask 115 (last digit 5 leaves world-write)",
			input: `umask 115`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1634",
					Message: "`umask 115` leaves world-write on new files — the \"other\" digit must be `2`/`3`/`6`/`7` to mask the write bit. Use `022` for public, `077` for secrets.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1634")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1635(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — mysql -p (prompts)",
			input:    `mysql -u root -p -h db.example.com`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — mysql --login-path",
			input:    `mysql --login-path=prod mydb`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — mysql -psecret",
			input: `mysql -u root -psecret -h db.example.com`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1635",
					Message: "`mysql -psecret` puts the MySQL password in argv. Use `-p` with no arg (prompt), `--login-path`, or a 0600 `~/.my.cnf`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — mysqldump -p$PW",
			input: `mysqldump -u root -p$PW mydb`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1635",
					Message: "`mysqldump -p$PW` puts the MySQL password in argv. Use `-p` with no arg (prompt), `--login-path`, or a 0600 `~/.my.cnf`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1635")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1636(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — virsh shutdown",
			input:    `virsh shutdown my-vm`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — virsh destroy --graceful",
			input:    `virsh destroy --graceful my-vm`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — virsh destroy my-vm",
			input: `virsh destroy my-vm`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1636",
					Message: "`virsh destroy` yanks power from the VM — filesystem corruption risk. Use `virsh shutdown` for graceful stop, or `virsh destroy --graceful` as a timed fallback.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — virsh destroy $DOM",
			input: `virsh destroy $DOM`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1636",
					Message: "`virsh destroy` yanks power from the VM — filesystem corruption risk. Use `virsh shutdown` for graceful stop, or `virsh destroy --graceful` as a timed fallback.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1636")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1637(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — unrelated command",
			input:    `export FOO=bar`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — readonly FOO=bar",
			input: `readonly FOO=bar`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1637",
					Message: "`readonly` works but `typeset -r NAME=value` is the Zsh-native form and composes with other typeset flags.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — readonly MAX_RETRIES=5",
			input: `readonly MAX_RETRIES=5`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1637",
					Message: "`readonly` works but `typeset -r NAME=value` is the Zsh-native form and composes with other typeset flags.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1637")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1638(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — non-secret build-arg",
			input:    `docker build --build-arg VERSION=1.0 -t app .`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — BuildKit secret",
			input:    `docker build --secret id=dbpass,src=/run/secrets/db -t app .`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker build --build-arg PASSWORD=s3cret",
			input: `docker build --build-arg PASSWORD=s3cret -t app .`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1638",
					Message: "`docker build --build-arg PASSWORD=s3cret` bakes the secret into the image layer metadata. Use `--secret id=NAME,src=PATH` (BuildKit) or a multi-stage build.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — podman build --build-arg API_KEY=$KEY",
			input: `podman build --build-arg API_KEY=$KEY -t app .`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1638",
					Message: "`podman build --build-arg API_KEY=$KEY` bakes the secret into the image layer metadata. Use `--secret id=NAME,src=PATH` (BuildKit) or a multi-stage build.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1638")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1639(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — non-credential header",
			input:    `curl -H "Content-Type: application/json" https://api`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — -H @file (read from file)",
			input:    `curl -H @/run/secrets/auth_header https://api`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — Authorization Bearer",
			input: `curl -H "Authorization: Bearer $TOKEN" https://api`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1639",
					Message: "`curl -H \"Authorization: Bearer $TOKEN\"` places the credential in argv — visible via `ps`. Use `-H @FILE` or `--config FILE` with 0600 perms.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — X-Api-Key",
			input: `curl -H "X-Api-Key: $KEY" https://api`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1639",
					Message: "`curl -H \"X-Api-Key: $KEY\"` places the credential in argv — visible via `ps`. Use `-H @FILE` or `--config FILE` with 0600 perms.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1639")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1640(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — Zsh (P) flag",
			input:    `echo "${(P)var}"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — plain expansion",
			input:    `echo "${var}"`,
			expected: []katas.Violation{},
		},
		{
			name:  `invalid — echo "${!var}"`,
			input: `echo "${!var}"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1640",
					Message: "`${!var}` Bash indirect — prefer Zsh `${(P)var}` for the same semantics with flag composability.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  `invalid — print "${!array[@]}"`,
			input: `print -r -- "${!array[@]}"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1640",
					Message: "`${!var}` Bash indirect — prefer Zsh `${(P)var}` for the same semantics with flag composability.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1640")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1641(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — --from-file",
			input:    `kubectl create secret generic mysec --from-file=password=/run/secrets/pw`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — --from-env-file",
			input:    `kubectl create secret generic mysec --from-env-file=/run/secrets/env`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — --from-literal=password=X",
			input: `kubectl create secret generic mysec --from-literal=password=hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1641",
					Message: "`kubectl create secret --from-literal=password=hunter2` puts the secret in argv — visible via `ps`. Use `--from-file=KEY=PATH` / `--from-env-file=PATH`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — --docker-password=$PW",
			input: `kubectl create secret docker-registry reg --docker-password=$PW --docker-username=u`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1641",
					Message: "`kubectl create secret --docker-password=$PW` puts the secret in argv — visible via `ps`. Use `--from-file=KEY=PATH` / `--from-env-file=PATH`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1641")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1642(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — tshark -w with -Z user",
			input:    `tshark -i eth0 -w /tmp/cap.pcap -Z analyst`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — tshark without -w (display only)",
			input:    `tshark -i any`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — tshark -w without -Z",
			input: `tshark -i any -w /tmp/cap.pcap`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1642",
					Message: "`tshark -w FILE` without `-Z USER` leaves the pcap root-owned. Add `-Z USER` to drop privileges for the capture.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — dumpcap -w without -Z",
			input: `dumpcap -i eth0 -w /tmp/cap.pcap`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1642",
					Message: "`dumpcap -w FILE` without `-Z USER` leaves the pcap root-owned. Add `-Z USER` to drop privileges for the capture.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1642")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1643(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — $(<file)",
			input:    `echo "$(<file)"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — bare cat file (logging, not capture)",
			input:    `cat /etc/os-release`,
			expected: []katas.Violation{},
		},
		{
			name:  `invalid — echo "$(cat /etc/hostname)"`,
			input: `echo "$(cat /etc/hostname)"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1643",
					Message: "`$(cat FILE)` forks cat just to read a file — use `$(<FILE)` (shell builtin, no fork).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  `invalid — print -r -- "$(cat file)"`,
			input: `print -r -- "$(cat file)"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1643",
					Message: "`$(cat FILE)` forks cat just to read a file — use `$(<FILE)` (shell builtin, no fork).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1643")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1644(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — unzip without -P (prompts)",
			input:    `unzip archive.zip`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — zip without password",
			input:    `zip archive.zip files/`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — unzip -P secret",
			input: `unzip -P s3cret archive.zip`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1644",
					Message: "`unzip -P` places the archive password in argv — visible via `ps`. Drop `-P` for interactive prompt, or switch to `7z -p` (reads from stdin) / `age` / `gpg` with keys in a protected file.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — zip -Psecret",
			input: `zip -Ps3cret archive.zip files/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1644",
					Message: "`zip -P` places the archive password in argv — visible via `ps`. Drop `-P` for interactive prompt, or switch to `7z -p` (reads from stdin) / `age` / `gpg` with keys in a protected file.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1644")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1645(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — source os-release",
			input:    `source /etc/os-release`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — lsb_release -rs",
			input: `lsb_release -rs`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1645",
					Message: "`lsb_release` needs an optional package. Use `source /etc/os-release` and read `$ID` / `$VERSION_ID` / `$PRETTY_NAME` instead — always present, no fork.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — lsb_release -a",
			input: `lsb_release -a`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1645",
					Message: "`lsb_release` needs an optional package. Use `source /etc/os-release` and read `$ID` / `$VERSION_ID` / `$PRETTY_NAME` instead — always present, no fork.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1645")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1646(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — btrfs scrub",
			input:    `btrfs scrub start /mnt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — btrfs check (read-only)",
			input:    `btrfs check $DEV`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — xfs_repair without -L",
			input:    `xfs_repair $DEV`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — btrfs check --repair",
			input: `btrfs check --repair $DEV`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1646",
					Message: "`btrfs check --repair` may worsen damage — try `btrfs scrub` and read-only `btrfs check` first, and snapshot the block device before running.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — xfs_repair -L",
			input: `xfs_repair -L $DEV`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1646",
					Message: "`xfs_repair -L` zeroes the log — uncommitted transactions are lost. Snapshot the block device first; mount read-only and read the log if possible.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1646")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1647(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — local file",
			input:    `kubectl apply -f ./manifest.yaml`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — stdin",
			input:    `kubectl apply -f -`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — kubectl apply -f https://example.com/m.yaml",
			input: `kubectl apply -f https://example.com/manifest.yaml`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1647",
					Message: "`kubectl apply -f https://example.com/manifest.yaml` applies a remote manifest — verify digest first. Download, check SHA256, then apply the local file.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — kubectl create -f http://insecure/m.yaml",
			input: `kubectl create -f http://insecure/m.yaml`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1647",
					Message: "`kubectl create -f http://insecure/m.yaml` applies a remote manifest — verify digest first. Download, check SHA256, then apply the local file.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1647")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1648(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — logrotate",
			input:    `logrotate -f /etc/logrotate.d/app`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — cp /dev/null to app tmp",
			input:    `cp /dev/null /tmp/marker`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — cp /dev/null /var/log/auth.log",
			input: `cp /dev/null /var/log/auth.log`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1648",
					Message: "`cp /dev/null /var/log/auth.log` wipes an audit log — use `logrotate -f` or `journalctl --vacuum-time=...` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — truncate -s 0 /var/log/wtmp",
			input: `truncate -s 0 /var/log/wtmp`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1648",
					Message: "`truncate -s 0 /var/log/wtmp` wipes an audit log — use `logrotate -f` or `journalctl --vacuum-time=...` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1648")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1649(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — -days 90 (Let's Encrypt style)",
			input:    `openssl req -x509 -days 90 -nodes`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — -days 398 (1-year max)",
			input:    `openssl req -x509 -days 398 -nodes`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — -days 3650 (10 years)",
			input: `openssl req -x509 -days 3650 -nodes -newkey rsa:2048`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1649",
					Message: "`openssl req -days 3650` issues a cert with a long validity. Keep leaf certs under 398 days and automate rotation.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — openssl x509 -days 1095 (3 years)",
			input: `openssl x509 -req -days 1095 -signkey key -in csr`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1649",
					Message: "`openssl x509 -days 1095` issues a cert with a long validity. Keep leaf certs under 398 days and automate rotation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1649")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1650(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — setopt NO_NOMATCH",
			input:    `setopt NO_NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — unsetopt BEEP",
			input:    `unsetopt BEEP`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — setopt RM_STAR_SILENT",
			input: `setopt RM_STAR_SILENT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1650",
					Message: "`setopt RM_STAR_SILENT` removes the `rm *` confirmation prompt — keep the default `RM_STAR_WAIT` so accidental deletions pause before they happen.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — unsetopt rmstarwait (lowercase, no underscore)",
			input: `unsetopt rmstarwait`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1650",
					Message: "`unsetopt rmstarwait` removes the `rm *` confirmation prompt — keep the default `RM_STAR_WAIT` so accidental deletions pause before they happen.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1650")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1651(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — localhost bind",
			input:    `docker run -p 127.0.0.1:8080:80 nginx`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — implicit port (not explicit 0.0.0.0)",
			input:    `docker run -p 8080:80 nginx`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — 0.0.0.0 explicit",
			input: `docker run -p 0.0.0.0:8080:80 nginx`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1651",
					Message: "`docker run -p 0.0.0.0:8080:80` publishes to every interface. Bind to `127.0.0.1:HOST:CONT` and put nginx / caddy in front for external access.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — [::] IPv6",
			input: `podman run -p [::]:8080:80 nginx`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1651",
					Message: "`podman run -p [::]:8080:80` publishes to every interface. Bind to `127.0.0.1:HOST:CONT` and put nginx / caddy in front for external access.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1651")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1652(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ssh without X11",
			input:    `ssh user@host`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ssh -X (untrusted)",
			input:    `ssh -X user@host firefox`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ssh -Y user@host",
			input: `ssh -Y user@host xclock`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1652",
					Message: "`ssh -Y` enables trusted X11 forwarding — remote clients get full access to the local X server. Use `-X` (untrusted) or drop X11 forwarding entirely.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ssh -i key -Y user@host",
			input: `ssh -i key -Y user@host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1652",
					Message: "`ssh -Y` enables trusted X11 forwarding — remote clients get full access to the local X server. Use `-X` (untrusted) or drop X11 forwarding entirely.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1652")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1653(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — $sysparams[pid]",
			input:    `echo "$sysparams[pid]"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — plain PID reference",
			input:    `echo "$PPID"`,
			expected: []katas.Violation{},
		},
		{
			name:  `invalid — echo "$BASHPID"`,
			input: `echo "$BASHPID"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1653",
					Message: "`$BASHPID` is Bash-only. Use `$sysparams[pid]` after `zmodload zsh/system`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  `invalid — print "sub=${BASHPID}"`,
			input: `print -r -- "sub=${BASHPID}"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1653",
					Message: "`$BASHPID` is Bash-only. Use `$sysparams[pid]` after `zmodload zsh/system`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1653")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1654(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — sysctl -p /etc/sysctl.conf",
			input:    `sysctl -p /etc/sysctl.conf`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — sysctl -p (default)",
			input:    `sysctl -p`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — sysctl -p /tmp/sysctl.conf",
			input: `sysctl -p /tmp/sysctl.conf`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1654",
					Message: "`sysctl -p /tmp/sysctl.conf` reads tunables from a world-traversable path — a concurrent local user can substitute the file. Keep configs under `/etc/sysctl.d/`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — sysctl -p /var/tmp/x",
			input: `sysctl -p /var/tmp/x.conf`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1654",
					Message: "`sysctl -p /var/tmp/x.conf` reads tunables from a world-traversable path — a concurrent local user can substitute the file. Keep configs under `/etc/sysctl.d/`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1654")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1655(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — read -k 1 (Zsh)",
			input:    `read -k 1 var`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — read -r line",
			input:    `read -r line`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — read -n 1 char",
			input: `read -n 1 char`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1655",
					Message: "`read -n N` is Bash syntax for \"read N characters\". Zsh's `-n` means \"drop trailing newline\" with no count. Use `read -k N var` in Zsh.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — read -n5 var",
			input: `read -n5 var`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1655",
					Message: "`read -n N` is Bash syntax for \"read N characters\". Zsh's `-n` means \"drop trailing newline\" with no count. Use `read -k N var` in Zsh.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1655")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1656(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — rsync with plain ssh",
			input:    `rsync -e ssh src host:dst`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — rsync without -e",
			input:    `rsync -a src dst`,
			expected: []katas.Violation{},
		},
		{
			name:  `invalid — rsync -e "ssh -o StrictHostKeyChecking=no"`,
			input: `rsync -e "ssh -o StrictHostKeyChecking=no" src host:dst`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1656",
					Message: "`rsync -e 'ssh -o StrictHostKeyChecking=no'` disables host-key verification — MITM risk. Pre-provision `known_hosts` and keep `StrictHostKeyChecking=yes`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  `invalid — rsync with UserKnownHostsFile=/dev/null`,
			input: `rsync -e "ssh -o UserKnownHostsFile=/dev/null" src host:dst`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1656",
					Message: "`rsync -e 'ssh -o StrictHostKeyChecking=no'` disables host-key verification — MITM risk. Pre-provision `known_hosts` and keep `StrictHostKeyChecking=yes`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1656")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1657(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — semanage permissive -d removes domain",
			input:    `semanage permissive -d httpd_t`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — semanage boolean -l listing",
			input:    `semanage boolean -l`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — semanage permissive -a httpd_t",
			input: `semanage permissive -a httpd_t`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1657",
					Message: "`semanage permissive -a` puts an SELinux domain in permissive mode — policy violations log but no longer block. Write a scoped allow rule with `audit2allow` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — semanage permissive --add sshd_t",
			input: `semanage permissive --add sshd_t`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1657",
					Message: "`semanage permissive -a` puts an SELinux domain in permissive mode — policy violations log but no longer block. Write a scoped allow rule with `audit2allow` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1657")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1658(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — curl -O without -J",
			input:    `curl -O https://example.com/file`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — curl -o with fixed name",
			input:    `curl -o out.bin https://example.com/file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — curl -OJ combined",
			input: `curl -OJ https://example.com/file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1658",
					Message: "`curl -OJ` saves the response under the name the server picks in `Content-Disposition` — path traversal is blocked but arbitrary same-dir overwrites are not. Pass `-o NAME` with a filename you control.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — curl -O -J split",
			input: `curl -O -J https://example.com/file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1658",
					Message: "`curl -OJ` saves the response under the name the server picks in `Content-Disposition` — path traversal is blocked but arbitrary same-dir overwrites are not. Pass `-o NAME` with a filename you control.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1658")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1659(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — fuser -k port target",
			input:    `fuser -k 8080/tcp`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — fuser without -k",
			input:    `fuser /var/log/syslog`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — fuser -k /mnt",
			input: `fuser -k /mnt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1659",
					Message: "`fuser -k /mnt` signals every process with a file open anywhere under the path — use PID / port targets or `systemctl stop` for services.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — fuser -kim /",
			input: `fuser -kim /`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1659",
					Message: "`fuser -k /` signals every process with a file open anywhere under the path — use PID / port targets or `systemctl stop` for services.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1659")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1660(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — printf %d without width",
			input:    `printf '%d' 5`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — printf %-20s left-aligned string (space fill)",
			input:    `printf '%-20s' "$name"`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — printf %05d",
			input: `printf '%05d' $n`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1660",
					Message: "`printf '%0Nd'` forks for zero-padding — prefer Zsh `${(l:N::0:)n}` parameter-expansion pad (same for `(r:N::0:)` on the right).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — printf %03d literal",
			input: `printf '%03d' 7`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1660",
					Message: "`printf '%0Nd'` forks for zero-padding — prefer Zsh `${(l:N::0:)n}` parameter-expansion pad (same for `(r:N::0:)` on the right).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1660")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1661(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — curl with real CA bundle",
			input:    `curl https://example.com --cacert /etc/ssl/certs/ca-certificates.crt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — curl without cacert",
			input:    `curl https://example.com -o out`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — curl URL --cacert /dev/null",
			input: `curl https://example.com --cacert /dev/null`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1661",
					Message: "`curl --cacert /dev/null` feeds curl an empty trust store — most TLS backends then accept any peer cert. Use a real bundle or `--pinnedpubkey`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — curl -s URL --capath /dev/null",
			input: `curl -s https://example.com --capath /dev/null`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1661",
					Message: "`curl --cacert /dev/null` feeds curl an empty trust store — most TLS backends then accept any peer cert. Use a real bundle or `--pinnedpubkey`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1661")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1662(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — pkexec direct command",
			input:    `pkexec /usr/bin/systemctl restart unit`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — pkexec apt install",
			input:    `pkexec /usr/bin/apt install foo`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — pkexec env DISPLAY=... cmd",
			input: `pkexec env DISPLAY=:0 /usr/bin/cmd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1662",
					Message: "`pkexec env VAR=VAL CMD` hands the root session a caller-controlled environment — use a polkit rule or `systemd-run --user` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — pkexec env PATH=/tmp cmd",
			input: `pkexec env PATH=/tmp /usr/bin/cmd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1662",
					Message: "`pkexec env VAR=VAL CMD` hands the root session a caller-controlled environment — use a polkit rule or `systemd-run --user` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1662")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1663(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — tune2fs -c 30 (reduced cadence)",
			input:    `tune2fs -c 30 $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — tune2fs -l (listing)",
			input:    `tune2fs -l $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — tune2fs -c 0",
			input: `tune2fs -c 0 $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1663",
					Message: "`tune2fs -c 0` disables periodic fsck on the filesystem — lower the cadence (e.g. `-c 30` / `-i 3m`) instead of turning it off.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — tune2fs -i 0",
			input: `tune2fs -i 0 $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1663",
					Message: "`tune2fs -i 0` disables periodic fsck on the filesystem — lower the cadence (e.g. `-c 30` / `-i 3m`) instead of turning it off.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1663")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1664(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — systemctl set-default multi-user.target",
			input:    `systemctl set-default multi-user.target`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — systemctl set-default graphical.target",
			input:    `systemctl set-default graphical.target`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — systemctl set-default rescue.target",
			input: `systemctl set-default rescue.target`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1664",
					Message: "`systemctl set-default rescue.target` makes every subsequent boot land in single-user mode — revert with `set-default multi-user.target` or `graphical.target`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — systemctl set-default emergency.target",
			input: `systemctl set-default emergency.target`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1664",
					Message: "`systemctl set-default emergency.target` makes every subsequent boot land in single-user mode — revert with `set-default multi-user.target` or `graphical.target`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1664")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1665(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — chrt -o 0 (SCHED_OTHER)",
			input:    `chrt -o 0 /usr/bin/cmd`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — chrt listing priority",
			input:    `chrt -p 1234`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — chrt -r 99 cmd",
			input: `chrt -r 99 /usr/bin/cmd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1665",
					Message: "`chrt -r` puts the child on a real-time scheduling class — a busy-loop or deadlock then starves kworker / sshd. Prefer `nice -n -5` or a systemd unit with `CPUWeight=`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — chrt -f 50 cmd",
			input: `chrt -f 50 /usr/bin/cmd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1665",
					Message: "`chrt -f` puts the child on a real-time scheduling class — a busy-loop or deadlock then starves kworker / sshd. Prefer `nice -n -5` or a systemd unit with `CPUWeight=`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1665")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1666(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — kubectl patch --type=strategic",
			input:    `kubectl patch deployment nginx --type=strategic`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — kubectl patch --type=merge",
			input:    `kubectl patch deployment nginx --type=merge`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — kubectl patch --type=json joined",
			input: `kubectl patch deployment nginx --type=json -p '[...]'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1666",
					Message: "`kubectl patch --type=json` applies a raw RFC-6902 patch that bypasses strategic-merge reconciliation — prefer `--type=strategic` and hold JSON patches behind code review.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — kubectl patch --type json split",
			input: `kubectl patch deployment nginx --type json -p '[...]'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1666",
					Message: "`kubectl patch --type=json` applies a raw RFC-6902 patch that bypasses strategic-merge reconciliation — prefer `--type=strategic` and hold JSON patches behind code review.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1666")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1667(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — openssl enc with -pbkdf2",
			input:    `openssl enc -aes-256-cbc -pbkdf2 -iter 100000 -in file -out file.enc`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — openssl req (different subcommand)",
			input:    `openssl req -new -key key.pem -out csr.pem`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — openssl enc -aes-256-cbc without pbkdf2",
			input: `openssl enc -aes-256-cbc -salt -in file -out file.enc`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1667",
					Message: "`openssl enc` without `-pbkdf2` uses single-round EVP_BytesToKey (MD5) — add `-pbkdf2 -iter 100000`, or prefer `age` / `gpg --symmetric`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — openssl enc -aes-256-gcm (no pbkdf2, no AEAD support)",
			input: `openssl enc -aes-256-gcm -in file -out file.enc`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1667",
					Message: "`openssl enc` without `-pbkdf2` uses single-round EVP_BytesToKey (MD5) — add `-pbkdf2 -iter 100000`, or prefer `age` / `gpg --symmetric`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1667")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1668(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — aws iam attach-user-policy ReadOnlyAccess",
			input:    `aws iam attach-user-policy --user-name foo --policy-arn arn:aws:iam::aws:policy/ReadOnlyAccess`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — aws iam create-access-key",
			input:    `aws iam create-access-key --user-name foo`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — attach-user-policy AdministratorAccess",
			input: `aws iam attach-user-policy --user-name foo --policy-arn arn:aws:iam::aws:policy/AdministratorAccess`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1668",
					Message: "`aws iam attach-user-policy ... arn:aws:iam::aws:policy/AdministratorAccess` grants sweeping admin — use a scoped inline policy (`put-user-policy`) or a customer-managed policy with the minimum `Action`/`Resource` set.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — attach-role-policy PowerUserAccess",
			input: `aws iam attach-role-policy --role-name r --policy-arn arn:aws:iam::aws:policy/PowerUserAccess`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1668",
					Message: "`aws iam attach-role-policy ... arn:aws:iam::aws:policy/PowerUserAccess` grants sweeping admin — use a scoped inline policy (`put-user-policy`) or a customer-managed policy with the minimum `Action`/`Resource` set.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1668")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1669(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — git gc default",
			input:    `git gc`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — git gc --auto",
			input:    `git gc --auto`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — git gc --prune=now",
			input: `git gc --prune=now --aggressive`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1669",
					Message: "`git gc --prune=now` bulldozes the reflog / prune recovery window — keep the default cadence unless you are actively purging leaked secrets, and mirror the dropped history off-box first.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — git reflog expire --expire=now --all",
			input: `git reflog expire --expire=now --all`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1669",
					Message: "`git reflog expire --expire=now` bulldozes the reflog / prune recovery window — keep the default cadence unless you are actively purging leaked secrets, and mirror the dropped history off-box first.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1669")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1670(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — setsebool -P httpd_can_network_connect on (not in dangerous list)",
			input:    `setsebool -P httpd_can_network_connect on`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — setsebool without -P (session only)",
			input:    `setsebool httpd_execmem 1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — setsebool -P httpd_execmem 1",
			input: `setsebool -P httpd_execmem 1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1670",
					Message: "`setsebool -P httpd_execmem 1` persistently relaxes SELinux memory-protection policy — fix the binary instead (`execstack -c`, relabel with `chcon`, or change the domain).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — setsebool -P allow_execstack on",
			input: `setsebool -P allow_execstack on`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1670",
					Message: "`setsebool -P allow_execstack on` persistently relaxes SELinux memory-protection policy — fix the binary instead (`execstack -c`, relabel with `chcon`, or change the domain).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1670")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1671(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — install -m 755",
			input:    `install -m 755 src /usr/local/bin/dst`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — mkdir -m 0755",
			input:    `mkdir -m 0755 /opt/dir`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — install -m 777",
			input: `install -m 777 src /usr/local/bin/dst`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1671",
					Message: "`install -m 777` creates a world-writable target — drop the world-write bit (e.g. `0755` / `0644`).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — mkdir -m 0666 (parser normalizes to 438)",
			input: `mkdir -m 0666 /shared`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1671",
					Message: "`mkdir -m 438` creates a world-writable target — drop the world-write bit (e.g. `0755` / `0644`).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1671")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1672(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — different command",
			input:    `semanage fcontext -a -t httpd_sys_content_t /var/www/app`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — chcon with no args",
			input:    `chcon`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — chcon -t unconfined_t path",
			input: `chcon -t unconfined_t /usr/local/bin/script`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1672",
					Message: "`chcon` writes an ephemeral SELinux label — `restorecon` / policy rebuild reverts it. Persist via `semanage fcontext -a -t TYPE 'REGEX'` + `restorecon`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — chcon -R -t bin_t dir",
			input: `chcon -R -t bin_t /srv/app`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1672",
					Message: "`chcon` writes an ephemeral SELinux label — `restorecon` / policy rebuild reverts it. Persist via `semanage fcontext -a -t TYPE 'REGEX'` + `restorecon`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1672")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1673(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — stty echo (restore)",
			input:    `stty echo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — stty raw",
			input:    `stty raw`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — stty -echo",
			input: `stty -echo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1673",
					Message: "`stty -echo` to mask password entry is fragile — a crash leaves the terminal echo-off. Use `read -s VAR` (Zsh / Bash 4+) instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1673")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1674(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker run default",
			input:    `docker run alpine`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — docker run --oom-score-adj 0",
			input:    `docker run --oom-score-adj 0 alpine`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker run --oom-kill-disable",
			input: `docker run --oom-kill-disable alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1674",
					Message: "`--oom-kill-disable` shifts OOM pressure onto the rest of the host — cap memory with `--memory=<limit>` instead of rigging the OOM score.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — podman run --oom-score-adj=-1000",
			input: `podman run --oom-score-adj=-1000 alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1674",
					Message: "`--oom-score-adj=-1000` shifts OOM pressure onto the rest of the host — cap memory with `--memory=<limit>` instead of rigging the OOM score.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1674")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1675(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — plain export VAR=value",
			input:    `export VAR=value`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — export multiple plain names",
			input:    `export PATH HOME`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — export -f function",
			input: `export -f my_func`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1675",
					Message: "`export -f` is Bash-only — use `typeset -fx` in Zsh.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — export -n VAR",
			input: `export -n VAR`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1675",
					Message: "`export -n` is Bash-only — use `typeset +x` in Zsh.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1675")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1676(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — helm rollback without --force",
			input:    `helm rollback myapp 2`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — helm history",
			input:    `helm history myapp`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — helm rollback --force",
			input: `helm rollback myapp 2 --force`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1676",
					Message: "`helm rollback --force` deletes and recreates unpatched resources — loses in-flight traffic and bypasses PodDisruptionBudget. Drop `--force` and gate the rollback via change review.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1676")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1677(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — trap cleanup EXIT",
			input:    `trap 'cleanup' EXIT`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — trap set -x on ERR (different signal)",
			input:    `trap 'set -x' ERR`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — trap 'set -x' DEBUG",
			input: `trap 'set -x' DEBUG`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1677",
					Message: "`trap 'set -x' DEBUG` keeps xtrace on after the first command — every subsequent argv (passwords, bearer tokens) lands in the log. Trace a narrow block with `set -x … set +x` or use `typeset -ft FUNC` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  `invalid — trap "set -o xtrace" DEBUG`,
			input: `trap "set -o xtrace" DEBUG`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1677",
					Message: "`trap 'set -x' DEBUG` keeps xtrace on after the first command — every subsequent argv (passwords, bearer tokens) lands in the log. Trace a narrow block with `set -x … set +x` or use `typeset -ft FUNC` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1677")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1678(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — borg init --encryption=repokey",
			input:    `borg init --encryption=repokey-blake2 /backup`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — borg list (different subcommand)",
			input:    `borg list /backup`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — borg init --encryption=none joined",
			input: `borg init --encryption=none /backup`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1678",
					Message: "`borg init --encryption=none` leaves archives unauthenticated and readable — use `--encryption=repokey-blake2` (or `keyfile-blake2`) and store the passphrase in `BORG_PASSPHRASE_FILE`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — borg init -e none",
			input: `borg init -e none /backup`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1678",
					Message: "`borg init --encryption=none` leaves archives unauthenticated and readable — use `--encryption=repokey-blake2` (or `keyfile-blake2`) and store the passphrase in `BORG_PASSPHRASE_FILE`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1678")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1679(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — gcloud add-iam-policy-binding roles/viewer",
			input:    `gcloud projects add-iam-policy-binding PROJ --member=user:a@ex.com --role=roles/viewer`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — gcloud iam service-accounts create",
			input:    `gcloud iam service-accounts create foo`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — add-iam-policy-binding roles/owner",
			input: `gcloud projects add-iam-policy-binding PROJ --member=user:a@ex.com --role=roles/owner`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1679",
					Message: "`gcloud ... add-iam-policy-binding --role=roles/owner` grants primitive / IAM-admin — use a predefined role with the minimum scope or a custom role, and apply admin changes via Terraform.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — add-iam-policy-binding roles/iam.serviceAccountTokenCreator",
			input: `gcloud projects add-iam-policy-binding PROJ --member=user:a@ex.com --role=roles/iam.serviceAccountTokenCreator`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1679",
					Message: "`gcloud ... add-iam-policy-binding --role=roles/iam.serviceAccountTokenCreator` grants primitive / IAM-admin — use a predefined role with the minimum scope or a custom role, and apply admin changes via Terraform.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1679")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1680(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — vault file under /etc/ansible",
			input:    `ansible-playbook site.yml --vault-password-file=/etc/ansible/vault.pass`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — no vault file",
			input:    `ansible-playbook site.yml`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — vault file under /tmp joined",
			input: `ansible-playbook site.yml --vault-password-file=/tmp/vault.pass`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1680",
					Message: "`ansible-playbook --vault-password-file` under `/tmp/` / `/var/tmp/` / `/dev/shm/` — world-traversable, any local user can race-read it. Store the key mode-0400 under `/etc/ansible/` or supply via a `vault-password-client` helper.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — vault file under /dev/shm split",
			input: `ansible-playbook site.yml --vault-password-file /dev/shm/vault`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1680",
					Message: "`ansible-playbook --vault-password-file` under `/tmp/` / `/var/tmp/` / `/dev/shm/` — world-traversable, any local user can race-read it. Store the key mode-0400 under `/etc/ansible/` or supply via a `vault-password-client` helper.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1680")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1681(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — tar -xf archive.tar",
			input:    `tar -xf archive.tar`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — tar -xvzf archive.tgz (no P)",
			input:    `tar -xvzf archive.tgz`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — tar -xPf archive.tar (short-flag cluster)",
			input: `tar -xPf archive.tar`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1681",
					Message: "`tar -xPf` keeps absolute paths during extraction — an untrusted archive can overwrite `/etc/cron.d`, `/usr/local/bin`, etc. Drop the flag and extract with `-C <scratch-dir>` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — tar xf archive.tar -P",
			input: `tar xf archive.tar -P`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1681",
					Message: "`tar -P` keeps absolute paths during extraction — an untrusted archive can overwrite `/etc/cron.d`, `/usr/local/bin`, etc. Drop the flag and extract with `-C <scratch-dir>` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1681")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1682(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — npm ci (no unsafe-perm)",
			input:    `npm ci`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — yarn install",
			input:    `yarn install`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — npm install --unsafe-perm",
			input: `npm install --unsafe-perm`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1682",
					Message: "`npm --unsafe-perm` keeps root for every lifecycle script — a compromised dep executes as root. Build in a dedicated builder container or run as a non-root user.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — npm install --unsafe-perm=true",
			input: `npm install --unsafe-perm=true`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1682",
					Message: "`npm --unsafe-perm=true` keeps root for every lifecycle script — a compromised dep executes as root. Build in a dedicated builder container or run as a non-root user.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1682")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1683(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — https registry",
			input:    `npm config set registry https://registry.npmjs.org/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — npm config set strict-ssl",
			input:    `npm config set strict-ssl false`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — npm config set registry http://...",
			input: `npm config set registry http://internal.example.com/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1683",
					Message: "`npm config set registry http://internal.example.com/` uses plaintext HTTP — any proxy / CDN can rewrite tarballs. Use `https://` and a custom CA via `NODE_EXTRA_CA_CERTS` if needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — yarn config set registry http://...",
			input: `yarn config set registry http://internal/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1683",
					Message: "`yarn config set registry http://internal/` uses plaintext HTTP — any proxy / CDN can rewrite tarballs. Use `https://` and a custom CA via `NODE_EXTRA_CA_CERTS` if needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1683")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1684(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — redis-cli PING",
			input:    `redis-cli PING`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — redis-cli -h host",
			input:    `redis-cli -h cache.example.com PING`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — redis-cli -a SECRET PING",
			input: `redis-cli -a SECRET PING`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1684",
					Message: "`redis-cli -a PASSWORD` leaks the password into `ps` / `/proc/PID/cmdline` — use `REDISCLI_AUTH` env var or `-askpass` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — redis-cli -aSECRET joined",
			input: `redis-cli -aSECRET PING`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1684",
					Message: "`redis-cli -a PASSWORD` leaks the password into `ps` / `/proc/PID/cmdline` — use `REDISCLI_AUTH` env var or `-askpass` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1684")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1685(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — sleep 30",
			input:    `sleep 30`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — sleep variable",
			input:    `sleep $timeout`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — sleep infinity",
			input: `sleep infinity`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1685",
					Message: "`sleep infinity` does not trap SIGTERM — the orchestrator hangs until SIGKILL. Use `exec tail -f /dev/null` or front with `tini` / `dumb-init`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1685")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1686(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — bare compinit",
			input:    `compinit -d $XDG_CACHE_HOME/zcompdump`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — different command",
			input:    `compaudit`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — compinit -C",
			input: `compinit -C`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1686",
					Message: "`compinit -C` (skip-security-check) loads `$fpath` files that are writable by others — any user on the host can inject shell code. Run `compaudit`, fix permissions, then `compinit` without the flag.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — compinit -u",
			input: `compinit -u`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1686",
					Message: "`compinit -u` (load-insecure-files) loads `$fpath` files that are writable by others — any user on the host can inject shell code. Run `compaudit`, fix permissions, then `compinit` without the flag.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1686")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1687(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — snap install strict",
			input:    `snap install hello-world`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — snap list",
			input:    `snap list`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — snap install --classic",
			input: `snap install code --classic`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1687",
					Message: "`snap install --classic` drops AppArmor / cgroup / seccomp sandbox — find a strict snap or a distro-package alternative, or document why this specific snap needs it.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — snap install --devmode",
			input: `snap install pkg --devmode`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1687",
					Message: "`snap install --devmode` logs confinement violations instead of blocking — find a strict snap or a distro-package alternative, or document why this specific snap needs it.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1687")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1688(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — aws s3 sync without --delete",
			input:    `aws s3 sync ./build s3://my-bucket/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — aws s3 cp",
			input:    `aws s3 cp file s3://bucket/`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — aws s3 sync --delete",
			input: `aws s3 sync ./build s3://my-bucket/ --delete`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1688",
					Message: "`aws s3 sync --delete` wipes DST objects that are missing from SRC — a mistyped SRC bulk-deletes the bucket. Scope to the prefix, dry-run first, or enable versioning + MFA-delete.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — aws s3 sync between buckets with --delete",
			input: `aws s3 sync s3://src/ s3://dst/ --delete`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1688",
					Message: "`aws s3 sync --delete` wipes DST objects that are missing from SRC — a mistyped SRC bulk-deletes the bucket. Scope to the prefix, dry-run first, or enable versioning + MFA-delete.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1688")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1689(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — borg delete without --force (prompts)",
			input:    `borg delete /backup::archive1`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — borg prune",
			input:    `borg prune --keep-last 7 /backup`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — borg delete --force",
			input: `borg delete --force /backup`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1689",
					Message: "`borg delete --force` skips confirmation and can nuke the whole repository on a typo — use `borg prune --keep-*` with a retention policy, or gate outright deletion behind a manual review.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1689")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1690(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — pip install normal package",
			input:    `pip install requests`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — pip install git+URL@commit-hash",
			input:    `pip install git+https://github.com/org/repo@abc1234`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — pip install git+URL@v1.2.3 tag",
			input:    `pip install git+https://github.com/org/repo@v1.2.3`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — pip install git+URL without ref",
			input: `pip install git+https://github.com/org/repo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1690",
					Message: "`pip install git+https://github.com/org/repo` tracks a moving git ref — pin to a commit SHA (`@abc1234…`) or signed tag (`@v1.2.3`), or use the PyPI release.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — pip install git+URL@main (moving branch)",
			input: `pip install git+https://github.com/org/repo@main`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1690",
					Message: "`pip install git+https://github.com/org/repo@main` tracks a moving git ref — pin to a commit SHA (`@abc1234…`) or signed tag (`@v1.2.3`), or use the PyPI release.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1690")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1691(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — rsync without --remove-source-files",
			input:    `rsync -av src/ dst/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — rsync --delete (different flag)",
			input:    `rsync -av --delete src/ dst/`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — rsync --remove-source-files (local)",
			input: `rsync -av --remove-source-files src/ dst/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1691",
					Message: "`rsync --remove-source-files` deletes SRC on optimistic per-file success — verify DST after the transfer and `rm` explicitly instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — rsync --remove-source-files (remote)",
			input: `rsync -a --remove-source-files host:src dst`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1691",
					Message: "`rsync --remove-source-files` deletes SRC on optimistic per-file success — verify DST after the transfer and `rm` explicitly instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1691")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1692(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — kexec -l load only",
			input:    `kexec -l /boot/vmlinuz --initrd=/boot/initrd.img`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — kexec -u unload",
			input:    `kexec -u`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — kexec -e",
			input: `kexec -e`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1692",
					Message: "`kexec -e` jumps to a preloaded kernel without firmware reboot — wtmp / auditd see nothing. Use `systemctl kexec` or a real `systemctl reboot` to keep the audit trail.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1692")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1693(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ionice -c 2 (best-effort)",
			input:    `ionice -c 2 -n 4 cmd`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ionice -c 3 (idle)",
			input:    `ionice -c 3 cmd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ionice -c 1 (real-time, split)",
			input: `ionice -c 1 -n 0 cmd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1693",
					Message: "`ionice -c 1` puts the child in the real-time I/O class — a long-running workload starves sshd / journald / the rest of the host. Stay on class 2.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ionice -c1 (real-time, joined)",
			input: `ionice -c1 cmd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1693",
					Message: "`ionice -c 1` puts the child in the real-time I/O class — a long-running workload starves sshd / journald / the rest of the host. Stay on class 2.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1693")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1694(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ssh without forwarding",
			input:    `ssh user@host`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ssh -J jump (ProxyJump)",
			input:    `ssh -J bastion user@target`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ssh -A",
			input: `ssh -A user@host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1694",
					Message: "`ssh -A` forwards the caller's `SSH_AUTH_SOCK` into the remote — any root on that host can reuse the keys. Use `ssh -J jumphost` instead, or a scoped key for the remote task.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ssh -o ForwardAgent=yes",
			input: `ssh -o ForwardAgent=yes user@host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1694",
					Message: "`ssh -o ForwardAgent=yes` forwards the caller's `SSH_AUTH_SOCK` into the remote — any root on that host can reuse the keys. Use `ssh -J jumphost` instead, or a scoped key for the remote task.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1694")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1695(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — terraform plan",
			input:    `terraform plan`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — terraform state mv (tracked rename)",
			input:    `terraform state mv old new`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — terraform state rm",
			input: `terraform state rm module.app.aws_instance.x`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1695",
					Message: "`terraform state rm` mutates shared state outside plan/apply — use `terraform import` or `apply -replace=ADDR` instead, and review / back up first.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — tofu state push",
			input: `tofu state push local.tfstate`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1695",
					Message: "`tofu state push` mutates shared state outside plan/apply — use `terraform import` or `apply -replace=ADDR` instead, and review / back up first.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1695")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1696(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — pnpm install --frozen-lockfile",
			input:    `pnpm install --frozen-lockfile`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — npm ci",
			input:    `npm ci`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — pnpm install --no-frozen-lockfile",
			input: `pnpm install --no-frozen-lockfile`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1696",
					Message: "`--no-frozen-lockfile` allows the lockfile to drift — the CI artifact no longer matches the reviewed dependency graph. Use `--frozen-lockfile` / `--immutable` in CI.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — yarn install --no-immutable",
			input: `yarn install --no-immutable`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1696",
					Message: "`--no-immutable` allows the lockfile to drift — the CI artifact no longer matches the reviewed dependency graph. Use `--frozen-lockfile` / `--immutable` in CI.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1696")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1697(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — cryptsetup open without --allow-discards",
			input:    `cryptsetup open $DISK data`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — cryptsetup luksClose",
			input:    `cryptsetup luksClose data`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — cryptsetup open --allow-discards",
			input: `cryptsetup open --allow-discards $DISK data`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1697",
					Message: "`cryptsetup --allow-discards` leaks free-sector layout to anyone with raw-device access — drop it if offline-disk inspection is in scope, or document the trade-off.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — cryptsetup luksOpen --allow-discards",
			input: `cryptsetup luksOpen --allow-discards $DISK data`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1697",
					Message: "`cryptsetup --allow-discards` leaks free-sector layout to anyone with raw-device access — drop it if offline-disk inspection is in scope, or document the trade-off.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1697")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1698(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — fail2ban-client status",
			input:    `fail2ban-client status`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — fail2ban-client set sshd unbanip scoped",
			input:    `fail2ban-client set sshd unbanip 1.2.3.4`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — fail2ban-client unban --all",
			input: `fail2ban-client unban --all`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1698",
					Message: "`fail2ban-client unban --all` wipes every active brute-force ban — attacker IPs regain access. Target individual IPs with `set <jail> unbanip <ip>` or reload a single jail.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — fail2ban-client stop",
			input: `fail2ban-client stop`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1698",
					Message: "`fail2ban-client stop` wipes every active brute-force ban — attacker IPs regain access. Target individual IPs with `set <jail> unbanip <ip>` or reload a single jail.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1698")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1699(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — kubectl drain with --ignore-daemonsets only",
			input:    `kubectl drain NODE --ignore-daemonsets`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — kubectl cordon (no drain)",
			input:    `kubectl cordon NODE`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — kubectl drain --delete-emptydir-data",
			input: `kubectl drain NODE --delete-emptydir-data --ignore-daemonsets`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1699",
					Message: "`kubectl drain --delete-emptydir-data` deletes `emptyDir` volumes along with the evicted pods — caches / WAL / scratch state are lost. Verify tolerance or migrate to a PersistentVolumeClaim first.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — kubectl drain --delete-local-data (deprecated alias)",
			input: `kubectl drain NODE --force --delete-local-data`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1699",
					Message: "`kubectl drain --delete-local-data` deletes `emptyDir` volumes along with the evicted pods — caches / WAL / scratch state are lost. Verify tolerance or migrate to a PersistentVolumeClaim first.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1699")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
