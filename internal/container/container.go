package container

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"path/filepath"
	"strconv"
)

/*
Get lightweight filesystem for Container:

wget https://dl-cdn.alpinelinux.org/alpine/v3.19/releases/x86_64/alpine-minirootfs-3.19.0-x86_64.tar.gz
mkdir -p giutainer-root
tar -xzf alpine-minirootfs-3.19.0-x86_64.tar.gz -C giutainer-root

*/

const (
	// Change according to personal container filesystem save folder
	pathfs = "/root/giutainer-root"
)

func Run() {
	args := append([]string{"child"}, os.Args[2:]...)

	cmd := exec.Command("/proc/self/exe", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS,
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Error Starting:", err)
		os.Exit(1)
	}
	if err := setupCgroup(cmd.Process.Pid); err != nil {
		fmt.Println("Error setting cgroup:", err)
		os.Exit(1)
	}
	if err := cmd.Wait(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func makedev(major, minor uint32) uint64 {
	return uint64(major)*256 + uint64(minor)
}

func Child() {
	fmt.Println("Giutainer - a minimal container runtime")

	if err := syscall.Chroot(pathfs); err != nil {
		fmt.Println("Error chroot:", err)
		os.Exit(1)
	}

	if err := os.Chdir("/"); err != nil {
		fmt.Println("Error chdir:", err)
		os.Exit(1)
	}

	// Mount /proc — makes ps, top, free work
	if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		fmt.Println("Error mounting proc:", err)
		os.Exit(1)
	}

	// Mount /dev as tmpfs
	if err := syscall.Mount("tmpfs", "/dev", "tmpfs", 0, ""); err != nil {
		fmt.Println("Error mounting dev:", err)
		os.Exit(1)
	}

	// Create essential device nodes
	devNull := syscall.Mknod("/dev/null", 0666|syscall.S_IFCHR, int(makedev(1, 3)))
	if devNull != nil {
		fmt.Println("Error creating /dev/null:", devNull)
	}

	devZero := syscall.Mknod("/dev/zero", 0666|syscall.S_IFCHR, int(makedev(1, 5)))
	if devZero != nil {
		fmt.Println("Error creating /dev/zero:", devZero)
	}

	devRandom := syscall.Mknod("/dev/random", 0666|syscall.S_IFCHR, int(makedev(1, 8)))
	if devRandom != nil {
		fmt.Println("Error creating /dev/random:", devRandom)
	}

	if err := syscall.Sethostname([]byte("giutainer")); err != nil {
		fmt.Println("Error setting hostname:", err)
		os.Exit(1)
	}

	// This makes it PID 1
	binary, err := exec.LookPath(os.Args[2])
	if err != nil {
		fmt.Println("Error finding command:", err)
		os.Exit(1)
	}

	if err := syscall.Exec(binary, os.Args[2:], os.Environ()); err != nil {
		fmt.Println("Error exec:", err)
		os.Exit(1)
	}
}

func setupCgroup(pid int) error {
	const cgroupPath = "/sys/fs/cgroup/giutainer"

	if err := os.MkdirAll(cgroupPath, 0755); err != nil {
		return fmt.Errorf("create cgroup: %w", err)
	}

	// Memory limit: 50MB
	if err := os.WriteFile(filepath.Join(cgroupPath, "memory.max"), []byte("52428800"), 0644); err != nil {
		return fmt.Errorf("set memory limit: %w", err)
	}

	// A period of 100000 microseconds (100ms) with a quota of 50000 means the process gets 50ms out of every 100ms. That’s 50% of one CPU core.
	// CPU limit: 50%
	if err := os.WriteFile(filepath.Join(cgroupPath, "cpu.max"), []byte("50000 100000"), 0644); err != nil {
		return fmt.Errorf("set cpu limit: %w", err)
	}

	if err := os.WriteFile(filepath.Join(cgroupPath, "cgroup.procs"), []byte(strconv.Itoa(pid)), 0644); err != nil {
		return fmt.Errorf("assign process: %w", err)
	}
	fmt.Printf("Cgroup: PID %d limited to 50MB memory, 50%% CPU\n", pid)
	return nil
}
