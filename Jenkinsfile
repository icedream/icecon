def binext(os) {
  switch(os) {
    case "windows":
      return ".exe"
    default:
      return ""
  }
}

def upx(file) {
  // Install upx
  tool "UPX v3.91"

  switch("${env.GOOS}.${env.GOARCH}") {
    case ~/linux\.(amd64|386)/:
    case ~/darwin\.(amd64|arm)/:
    case ~/windows\..+/:
    case ~/.+bsd\.386/:
      if (env.GOOS == "linux") {
        sh "GOOS= GOARCH= go get -v github.com/pwaller/goupx"
        sh "goupx --no-upx \"$file\""
      }
      sh "upx --best --ultra-brute \"$file\""
      break
    default:
      echo "Skipping UPX compression as it is not supported for $goos/$goarch."
      break
  }
}

def withGoEnv(os, arch, f) {
  // Install go
  env.GOROOT = tool "Go 1.7"
  env.GOPATH = env.WORKSPACE

  switch(arch) {
    case "x64":
      arch = "amd64"
      break
    case "x86":
      arch = "386"
      break
    case "armv5":
    case "armv6":
    case "armv7":
      arch = "arm"
      break
    case "armv8":
      arch = "arm64"
      break
  }

  withEnv(["CGO_ENABLED=1", "GOOS=${os}", "GOARCH=${arch}"]) {
    switch(arch) {
      case "armv5":
        withEnv("GOARM=5", f)
        break
      case "armv6":
        withEnv("GOARM=6", f)
        break
      case "armv7":
        withEnv("GOARM=7", f)
        break
      default:
        f()
        break
    }
  }
}

def build(os, arch) {
  node("docker && linux && amd64") {
    checkout scm
    docker.image("dockcross/${os}-${arch}").inside {
      withGoEnv(os, arch) {
        def binfilename = "icecon_${env.GOOS}_${env.GOARCH}${binext os}"
        sh "GOOS= GOARCH= go get -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo"
        sh "go generate -v ./..."
        sh "go get -v -d ./..."
        sh "go build -o ${binfilename}"
        upx binfilename
        archive "${binfilename}"
      }
    }
  }
}

parallel (
  windows_x64: { build("windows", "x64") },
  windows_x86: { build("windows", "x86") },
  linux_x64: { build("linux", "x64") },
  linux_x86: { build("linux", "x86") },
  linux_arm: { build("linux", "armv5") }
)