require 'formula'

class Butler < Formula
  url 'https://raw.github.com/fd/butler/master/homebrew/butler-0.0.2.tar.gz'
  homepage 'http://github.com/fd/butler'
  md5 'b2b54121ff705b502ca25e01df98bc35'

  skip_clean ['bin']

  def install
    bin.install "butler"
  end
end
