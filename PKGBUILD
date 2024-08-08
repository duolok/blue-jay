# Maintainer: duolok
pkgname=blue-jay
pkgver=0.0.1
pkgrel=1
pkgdesc="Minimalistic TUI for video game price tracking"
arch=('x86_64')
url="https://github.com/duolok/blue-jay"
license=('MIT')
depends=('go')  
makedepends=('git')  
source=("$pkgname-$pkgver.tar.gz::https://github.com/duolok/blue-jay/archive/refs/tags/v$pkgver.tar.gz")
sha256sums=('SKIP')  

build() {
  cd "$srcdir/$pkgname-$pkgver"
  go build -o blue-jay
}

package() {
  cd "$srcdir/$pkgname-$pkgver"
  install -Dm755 blue-jay "$pkgdir/usr/bin/blue-jay"
}

check() {
  cd "$srcdir/$pkgname-$pkgver"
  go test ./...
}

