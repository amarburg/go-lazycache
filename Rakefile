
task :default => :test

task :build do
  sh( *%w( go build ))
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

  task :short => :build do
      sh "go test -v"
  end

  task :integration => :build do
      sh "go test -v -tags integration"
  end

  task :redis => :build do
      sh "go test -v -tags redis"
  end

end


namespace :docker do
  task :build do
    sh "docker build -f Dockerfile --tag lazycache:dev ."
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
