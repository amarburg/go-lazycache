
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

task :test => "test:short"

namespace :test do
  task :all => ["test:integration","test:redis"]

  task :short => :build do
      sh *%w( go test -v  )
  end

  task :integration => :build do
      sh *%w( go test -v -tags integration )
  end

  task :redis => :build do
      sh *%w( go test -v -tags redis  )
  end

end


namespace :wercker do

  task :build do
    sh *%w( wercker --verbose build --git-domain github.com --git-owner=amarburg --git-repository=go-lazycache )
  end
end
