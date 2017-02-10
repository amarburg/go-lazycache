

task :build do
  sh *%w( go build )
end

task :test => :build do
    sh *%w( go test -tags integration )
end
