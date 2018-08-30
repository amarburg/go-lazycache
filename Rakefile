
task :default => :test

task :easyjson do
  sh "go get -u github.com/mailru/easyjson/..."
end

task :build => :easyjson do
  sh "go get -v"
  sh "easyjson -all moov_handler.go"
  sh "go build"
end

task :gofmt do
  sh "gofmt -s -w ."
end

task :lint do
  sh "golint ."
end

task :test => "test:short"

namespace :test do
  task :all => ["test:integration","test:redis"]

  task :deps do
    sh "go get -t"
  end

  task :short => [:build, "test:deps"] do
      sh "go test -v"
  end

  task :integration => [:build, "test:deps"] do
      sh "go test -v -tags integration"
  end

  task :redis => [:build, "test:deps"] do
      sh "go test -v -tags redis"
  end

end


namespace :docker do
  task :build do
    sh "docker build --file deploy/docker/Dockerfile --tag lazycache:dev ."
  end

  task :run do
    sh "docker run --rm --publish 8080:8080 lazycache:dev"
  end
end


namespace :wercker do

  desc "Build Wecker locally using wercker CLI"
  task :build do
    sh "wercker --verbose build --git-domain github.com --git-owner=amarburg --git-repository=go-lazycache"
  end
end
