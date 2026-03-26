# 🎵 MusicLe CLI Music Player

<div align="center">

**Terminal üzerinden müzik dinlemenizi sağlayan, hızlı ve şık bir müzik çalar.**

</div>

## 🌟 Hakkında

MusicLe, terminal üzerinden müzik dinleme deneyimini yeniden tanımlayan modern bir CLI müzik çalar'dır. Spotify-inspired arayüzü, hızlı performansı ve zengin özellikleri ile müzik koleksiyonunuzu terminalden yönetmenin en elegant yoludur.

### 🎯 Neden MusicLe?

- **🚀 Hızlı**: Go ile yazıldı, minimum kaynak kullanımı
- **🎨 Şık**: Spotify-inspired modern terminal arayüzü
- **🔥 Özellik Zengin**: Playlist yönetimi, Spotify entegrasyonu, yerel müzik desteği
- **🌡️ Hafif**: Sadece birkaç MB, anında başlangıç
- **🔧 Esnek**: Windows, macOS, Linux desteği

---

## 🚀 Hızlı Başlangıç

### Ön Gereksinimler

- **Go 1.26.1** veya üzeri
- **Git**
- **Terminal** (ANSI desteği olan)

### 30 Saniyede Kurulum

```bash
# 1. Repoyu klonla
git clone https://github.com/alperencolgecen/musicle-cli.git
cd musicle-cli

# 2. Derle ve çalıştır
go build -o musicle ./main.go
./musicle
```

İşte bu kadar! 🎉 MusicLe artık çalışıyor ve müzik dinlemeye hazır!

---

## 📖 Detaylı Kurulum

### Windows

```powershell
# Go yüklü değilse
winget install GoLang.Go

# Repoyu klonla
git clone https://github.com/alperencolgecen/musicle-cli.git
cd musicle-cli

# Derle
go build -o musicle.exe ./main.go

# Çalıştır
.\musicle.exe
```

### macOS

```bash
# Homebrew ile Go
brew install go

# Repoyu klonla
git clone https://github.com/alperencolgecen/musicle-cli.git
cd musicle-cli

# Derle ve çalıştır
go build -o musicle ./main.go
./musicle
```

### Linux

```bash
# Go yükleme (Ubuntu/Debian)
sudo apt update
sudo apt install golang-go

# Repoyu klonla
git clone https://github.com/alperencolgecen/musicle-cli.git
cd musicle-cli

# Derle ve çalıştır
go build -o musicle ./main.go
./musicle
```

---

## 🎮 Kullanım Rehberi

### İlk Kurulum

MusicLe'yi ilk çalıştırdığınızda 4 adımlık bir kurulum sihirbazı sizi karşılar:

1. **📁 Müzik Dizini**: Müzik dosyalarınızın saklanacağı konumu seçin
2. **🌐 Dil**: Arayüz dilini seçin (Türkçe/İngilizce)
3. **👤 Profil**: Profil adınızı ve görünen adınızı belirleyin
4. **📋 Playlist**: İlk playlist'inizi oluşturun

### Ana Arayüz

MusicLe'nin ana arayüzü 4 bölümden oluşur:

```
┌──────────────────────────────────────────────────────────────────────────┐
│  MusicLe          [Home]  [Settings]                                     │  ← Header
├─────────────────┬────────────────────────────────────────────────────────┤
│                 │  [Playlist ComboBox]                                   │
│  MUSIC DOWNLOAD │  [Playlist Art 70%] [Name] [Bio]                      │
│                 │  [🔒 Encrypt] [🔀 Shuffle] [▶ Play] [⬇ Download]      │
│  [Spotify URL ] │  ─────────────────────────────────────────────────    │
│  [YouTube URL ] │  #   Art   Title/Artist        Date Added   Duration  │
│  [+Local Music] │  ─────────────────────────────────────────────────    │
│  [Playlist ▾  ] │  1   🎵   Song Title           2024-01-01   03:45    │
│                 │       Artist Name                                      │
│  (1/4 width)    │  2   🎵   Song Title           2024-01-02   04:12    │
│                 │       Artist Name                                      │
│                 │  3   🎵   ...                  ...          ...       │
│                 │  4   🎵   ...                  ...          ...       │
│                 │  5   🎵   ...                  ...          ...       │
│                 │  6   🎵   ...                  ...          ...       │
│                 │                          (3/4 width — Spotify style)   │
├─────────────────┴────────────────────────────────────────────────────────┤
│ [AlbumArt] Song Title      ──●───────────────── 01:23 / 03:45   🔊████░ │
│            Artist Name                                                   │
└──────────────────────────────────────────────────────────────────────────┘
```

### Klavye Kısayolları

| Tuş | İşlev | Açıklama |
|-----|-------|----------|
| `Space` | ⏯️ Play/Pause | Çal/Durdur |
| `→` | ⏩ 5sn İleri | 5 saniye ileri sar |
| `←` | ⏪ 5sn Geri | 5 saniye geri sar |
| `↑` | 🔊 Ses Artır | Ses seviyesini artır |
| `↓` | 🔉 Ses Azalt | Ses seviyesini azalt |
| `Tab` | 🔄 Alan Değiştir | Bir sonraki alana geç |
| `Esc` | ❌ Çıkış | Uygulamadan çıkış |

### Müzik Ekleme

#### Spotify'dan Ekleme
1. Spotify URL'sini kopyala (örn: `https://open.spotify.com/track/...`)
2. Sidebar'daki Spotify alanına yapıştır
3. Enter'a bas

#### YouTube'dan Ekleme
1. YouTube Music URL'sini kopyala
2. Sidebar'daki YouTube alanına yapıştır
3. Enter'a bas

#### Yerel Müzik Ekleme
1. Sidebar'daki "+ Add Local Music" butonuna tıkla
2. Müzik dosyasının yolunu gir
3. Playlist seç ve ekle

---

## ⚙️ Özellikler

### 🎵 Müzik Yönetimi
- **📥 Çoklu Kaynak**: Spotify, YouTube Music, yerel dosyalar
- **📋 Playlist Yönetimi**: Oluştur, düzenle, sil
- **🔍 Arama**: Hızlı müzik arama
- **🏷️ Etiketleme**: Müzikleri kategorize et

### 🎨 Arayüz Özellikleri
- **🌙 Temalar**: Açık/Koyu tema desteği
- **📱 Responsive**: Farklı terminal boyutlarına uyum
- **🎨 Renkler**: Spotify-inspired renk paleti
- **⚡ Animasyonlar**: Smooth geçişler ve animasyonlar

### 🔧 Teknik Özellikler
- **⚡ Performans**: Go ile yüksek performans
- **🔒 Güvenlik**: Güvenli API entegrasyonu
- **📊 İstatistikler**: Çalma listeleri, dinleme süreleri
- **💾 Veri Yedekleme**: Profil ve playlist yedekleme

---

## 🏗️ Proje Yapısı

```
musicle-cli/
├── 📁 main.go              # Ana uygulama giriş noktası
├── 📁 internal/             # İç modüller
│   ├── 📁 bridge/           # UI ve motor arası köprü
│   ├── 📁 fs/               # Dosya sistemi işlemleri
│   ├── 📁 theme/            # UI temaları
│   └── 📁 ui/               # Kullanıcı arayüzü bileşenleri
│       ├── 📄 dashboard.go
│       ├── 📄 sidebar.go
│       ├── 📄 player_bar.go
│       └── 📄 ...
├── 📁 engine/               # Müzik çalar motoru
├── 📁 .github/              # GitHub şablonları
├── 📄 go.mod                # Go modül dosyası
├── 📄 README.md             # Bu dosya
├── 📄 CONTRIBUTING.md       # Katkı rehberi
├── 📄 CODE_OF_CONDUCT.md    # Davranış kuralları
└── 📄 SECURITY.md           # Güvenlik politikası
```

---

## 🔧 Geliştirme

### Geliştirme Ortamı Kurulumu

```bash
# Repoyu klonla
git clone https://github.com/alperencolgecen/musicle-cli.git
cd musicle-cli

# Modülleri yükle
go mod tidy

# Geliştirme için derle
go build -o musicle ./main.go

# Testleri çalıştır
go test ./...

# Linter çalıştır
golangci-lint run
```

### Kod Standartları

- **Go fmt**: Kod formatlaması için `go fmt`
- **Go vet**: Statik analiz için `go vet`
- **Test coverage**: Minimum %80 coverage
- **Documentation**: Public fonksiyonlara dokümantasyon

### Katkı Süreci

1. **🍴 Fork**: Repoyu fork et
2. **🌿 Branch**: Yeni branch oluştur (`git checkout -b feature/amazing-feature`)
3. **💻 Code**: Değişiklikleri yap
4. **🧪 Test**: Testleri çalıştır
5. **📤 PR**: Pull request gönder

---

## 📊 İstatistikler

<div align="center">

![GitHub stars](https://img.shields.io/github/stars/alperencolgecen/musicle?style=social)
![GitHub forks](https://img.shields.io/github/forks/alperencolgecen/musicle?style=social)
![GitHub issues](https://img.shields.io/github/issues/alperencolgecen/musicle)
![GitHub pull requests](https://img.shields.io/github/issues-pr/alperencolgecen/musicle)

</div>

### 📈 Proje Metrikleri
- **📝 Kod Satırı**: 5,000+ satır Go kodu
- **🧪 Test Coverage**: %85+
- **📦 Bağımlılıklar**: Minimum, sadece gerekli kütüphaneler
- **🚀 Performans**: <100ms başlangıç süresi

---

## 🤝 Katkı

MusicLe'ye katkıda bulunmak ister misiniz? Harika! İşte nasıl yapacağınız:

### 🎯 Nasıl Katkıda Bulunulur?

1. **🐛 Hata Raporları**: [Issue aç](https://github.com/alperencolgecen/musicle-cli/issues/new?template=bug_report.md)
2. **✨ Özellik İstekleri**: [Feature request](https://github.com/alperencolgecen/musicle-cli/issues/new?template=feature_request.md)
3. **💻 Kod Katkısı**: [CONTRIBUTING.md](CONTRIBUTING.md) rehberini takip et
4. **📖 Dokümantasyon**: Dokümantasyon geliştirme
5. **🌍 Çeviri**: Farklı dillere çeviri yap

## 🔒 Güvenlik

Güvenlik bizim için önemlidir. Güvenlik sorunlarını bildirmek için:

- **📧 Email**: alperencolgecen@gmail.com
- **🔐 Private**: Güvenlik sorunlarını public olarak bildirmeyin
- **⏰ Response**: 24 saat içinde yanıt

Detaylı bilgi için [SECURITY.md](SECURITY.md) dosyasını inceleyin.

---

## 📜 Lisans

Bu proje **Apache License 2.0** ile korunmaktadır.

### ✅ Yapabilirsiniz:
- ✅ Projeyi istediğiniz gibi kullanabilirsiniz
- ✅ Değişiklik yapabilir ve dağıtabilirsiniz
- ✅ Ticari amaçla kullanabilirsiniz
- ✅ Fork edebilir ve kendi projenizi yapabilirsiniz

### ❅ Yapamazsınız:
- ❌ Lisans ve telif hakkı bildirimlerini kaldıramazsınız
- ❌ Projeyi sahiplenemezsiniz
- ❌ Sorumluluk reddini değiştiremezsiniz

---

## 🌟 Gelecek Planları

### 🚀 Yakında Gelecek Özellikler

- [ ] **🎨 Tema Düzenleyici**: Özel tema oluşturma
- [ ] **📱 Mobil Uygulama**: iOS/Android companion app
- [ ] **🔌 Eklenti Sistemi**: Özelleştirilebilir eklentiler
- [ ] **🌐 Web Arayüzü**: Browser tabanlı arayüz
- [ ] **📊 Analitikler**: Detaylı dinleme istatistikleri
- [ ] **🔄 Sync**: Çoklu cihaz senkronizasyonu

### 🎯 Uzun Vadeli Hedefler

- **🏆 Lider Olmak**: En iyi CLI müzik çalar olmak
- **🌍 Global**: Uluslararası kullanıcı kitlesi
- **🔧 Entegrasyonlar**: Daha fazla servis entegrasyonu
- **📚 Eğitim**: Terminal tabanlı uygulamalar için eğitim

---

## 💬 İletişim

### 📧 İletişim Bilgileri

- **👨‍💻 Geliştirici**: Alperen Çölgeçen
- **📧 Email**: alperencolgecen@gmail.com
- **🐙 GitHub**: [@alperencolgecen](https://github.com/alperencolgecen)

### 🗨️ Topluluk

- **💬 GitHub Discussions**: Sorular ve tartışmalar
- **🐛 Issues**: Hata raporları ve özellik istekleri

---

## 🙏 Teşekkürler

Bu projeyi mümkün kılan herkese teşekkür ederiz:

- **🎵 Spotify**: İlham veren arayüz tasarımı
- **🔧 Go Team**: Harika programlama dili
- **🌟 tview**: Terminal UI kütüphanesi
- **👥 Katkıcular**: Değerli katkıları için
- **🎧 Kullanıcılar**: Destek ve feedback için

---

## 📄 Sürüm Notları

### v1.0.0 (Güncel)
- ✨ İlk sürüm
- 🎵 Spotify ve YouTube Music entegrasyonu
- 📋 Playlist yönetimi
- 🎨 Modern terminal arayüzü
- 🔧 Windows, macOS, Linux desteği

---

<div align="center">

---

**🎵 Terminalden müzik dinlemenin en şık yolu!**

Prepared by Alperen Çölgeçen

</div>
