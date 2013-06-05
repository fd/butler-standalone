require 'formula'

class Butler < Formula
  url 'https://github.com/fd/butler-standalone/raw/master/homebrew/butler-0.0.2.tar.gz'
  homepage 'https://github.com/fd/butler-standalone'
  sha256 '2a277342e6e047e73111d04a9f4e7a8a14948a65348cd09351360abee5b93e19'

  skip_clean ['bin']

  def install
    bin.install "butler"
  end
end
