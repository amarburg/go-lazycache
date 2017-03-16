
task :default => :test

task :build do
  sh *%w( go build )
end

task :gofmt do
  sh "gofmt -s -w ."
end

task :lint do
  sh "golint ."
end

task :test => :build do
    sh *%w( go test -v -tags integration )
end


namespace :wercker do

  task :build do
    sh *%w( wercker --verbose build --git-domain github.com --git-owner=amarburg --git-repository=go-lazycache )
  end
end
