require 'formula'

class Butler < Formula
  url 'https://raw.github.com/fd/butler/master/homebrew/butler-0.0.1.tar.gz'
  homepage 'http://github.com/fd/butler'
  md5 'e3dafa46a6bc08d056e1ee1a831a0e5f'

  skip_clean ['bin']

  def install
    bin.install "butler"
  end
end
