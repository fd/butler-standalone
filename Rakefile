
def with_env(env)
  _env = {}

  env.each do |k,v|
    _env[k] = ENV[k]
    ENV[k] = v
  end

  yield

ensure
  _env.each do |k,v|
    if v
      ENV[k] = v
    else
      ENV.delete(k)
    end
  end
end

GO_ROOT      = `go env GOROOT`.strip
GO_PATH      = File.expand_path('..', __FILE__)
GO_HOST_OS   = `go env GOHOSTOS`.strip
GO_HOST_ARCH = `go env GOHOSTARCH`.strip

OS_TARGETS   = %w( darwin linux freebsd openbsd windows )
ARCH_TARGETS = %w( 386 amd64 )

DEPENDENCIES = %w(
  github.com/nu7hatch/gopqueue
  github.com/dgryski/dgobloom
)

PACKAGES = %w(
  github.com/nu7hatch/gopqueue
  github.com/dgryski/dgobloom
)

TOOLS = {
  "butler" => "butler"
}

LIBS = Hash.new { |h,k| h[k] = [] }
BINS = Hash.new { |h,k| h[k] = [] }
PKGS = Hash.new { |h,k| h[k] = [] }

directory "dist"
OS_TARGETS.each do |os|
  ARCH_TARGETS.each do |arch|
    directory "bin/#{os}_#{arch}"
  end
end

PACKAGES.each do |package|
  OS_TARGETS.each do |os|
    ARCH_TARGETS.each do |arch|

      LIBS["#{os}_#{arch}"] << "pkg/#{os}_#{arch}/#{package}.a"
      deps = FileList["src/#{package}/*.go"]

      file "pkg/#{os}_#{arch}/#{package}.a" => deps do
        with_env "GOPATH" => GO_PATH, "GOOS" => os, "GOARCH" => arch, "CGO_ENABLED" => "0" do
          sh "go install #{package}"
        end
      end

    end
  end
end

TOOLS.each do |package, tool|
  OS_TARGETS.each do |os|
    ARCH_TARGETS.each do |arch|

      BINS["#{os}_#{arch}"] << "bin/#{os}_#{arch}/#{tool}"
      deps  = FileList["src/#{package}/*.go"]
      deps += ["bin/#{os}_#{arch}"]
      deps += LIBS["#{os}_#{arch}"]

      file "bin/#{os}_#{arch}/#{tool}" => deps do
        with_env "GOPATH" => GO_PATH, "GOOS" => os, "GOARCH" => arch, "CGO_ENABLED" => "0" do
          sh "go build -o bin/#{os}_#{arch}/#{tool} src/#{package}/*.go"
        end
      end

      PKGS["#{os}_#{arch}"] << "dist/#{tool}-#{os}_#{arch}.tar.gz"
      deps  = ["bin/#{os}_#{arch}/#{tool}"]
      deps += ["dist"]

      file "dist/#{tool}-#{os}_#{arch}.tar.gz" => deps do
        dst = File.expand_path("dist/#{tool}-#{os}_#{arch}.tar.gz")
        Dir.chdir("bin/#{os}_#{arch}") do
          sh "tar czf #{dst} #{tool}"
        end
      end

    end
  end
end

task :clean do
  sh "rm -rf bin"
  sh "rm -rf pkg"
  sh "rm -rf dist"
end

task :update do
  DEPENDENCIES.each do |dep|
    with_env "GOPATH" => GO_PATH do
      sh "go get -d -u #{dep}"
    end
  end
end

task :build => BINS.values.flatten

task :package => PKGS.values.flatten

task :local => BINS["#{GO_HOST_OS}_#{GO_HOST_ARCH}"]

task :setup do
  Dir.chdir(GO_ROOT + '/src') do

    %w( 8 6 ).each do |c|
      %w( a c g l ).each do |t|
        sh "go tool dist install -v cmd/#{c}#{t}"
      end
    end

    OS_TARGETS.each do |os|
      ARCH_TARGETS.each do |arch|
        with_env "GOOS" => os, "GOARCH" => arch do
          sh "go tool dist install -v pkg/runtime"
          sh "go install -v -a std || true"
        end
      end
    end

  end
end

task :default => :build
