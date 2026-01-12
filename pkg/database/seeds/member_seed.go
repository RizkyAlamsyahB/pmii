package seeds

import (
	"encoding/json"
	"log"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
)

// MemberSeedData contains member seeding data
type MemberSeedData struct {
	FullName    string
	Position    string
	Department  string
	PhotoFile   string
	SocialLinks string // JSON string
}

// SeedMembers seeds member data
func (s *Seeder) SeedMembers() error {
	logSeederStart("Members")

	// Clear existing members and reset sequence for clean re-seeding
	s.db.Exec("TRUNCATE TABLE members RESTART IDENTITY CASCADE")

	members := getMembersData()
	successCount := 0

	for i, m := range members {
		// Upload photo
		photoPath := getFilePath(s.seedsPath, "members", m.PhotoFile)
		photoURL, err := uploadFile(s.uploader, photoPath, "members")
		if err != nil {
			log.Printf("⚠️ Warning: Failed to upload photo for %s: %v", m.FullName, err)
			continue
		}

		// Parse social links
		var socialLinks map[string]any
		if err := json.Unmarshal([]byte(m.SocialLinks), &socialLinks); err != nil {
			log.Printf("⚠️ Warning: Failed to parse social links for %s: %v", m.FullName, err)
			socialLinks = make(map[string]any)
		}

		member := domain.Member{
			FullName:    m.FullName,
			Position:    m.Position,
			Department:  domain.MemberDepartment(m.Department),
			PhotoURI:    &photoURL,
			SocialLinks: socialLinks,
			IsActive:    true,
		}

		if err := s.db.Create(&member).Error; err != nil {
			log.Printf("⚠️ Warning: Failed to create member %s: %v", m.FullName, err)
			continue
		}

		successCount++
		logSeederProgress(i+1, len(members), m.FullName)
	}

	logSeederResult("Members", successCount, len(members))
	return nil
}

// getMembersData returns all member seed data
func getMembersData() []MemberSeedData {
	return []MemberSeedData{
		// Pengurus Harian
		{"M. Shofiyullah Cokro", "Ketua Umum", "pengurus_harian", "1765939378.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/"}`},
		{"M. Irkham Thamrin", "Sekretaris Jenderal", "pengurus_harian", "1765939447.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/"}`},
		{"Sainuddin", "Bendahara Umum", "pengurus_harian", "1765939452.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/"}`},

		// Kabid
		{"Acep Jamaluddin", "Ketua Bidang Kaderisasi Nasional", "kabid", "1765939456.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/cepjam"}`},
		{"Moh Lutfi", "Ketua Bidang Penataan Aparatur", "kabid", "1765939462.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/lutfi164"}`},
		{"M. Muham Tashir", "Ketua Bidang Okp Kemahasiswaan Lsm Dan Ormas", "kabid", "1765939467.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/"}`},
		{"Syamsuddin", "Ketua Bidang Hubungan Agama Dan Hubungan Umat Beragama", "kabid", "1765939472.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/syamsuddin_jm"}`},
		{"Sibly Adam Firnanda", "Ketua Bidang Hubungan Internasional", "kabid", "1765939477.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/syibly29","linkedin":"https://linkedin.com/in/syibly-adam-firmanda"}`},
		{"Aprilana Eka Dani", "Ketua Bidang Pendidikan, Riset, Ilmu Pengetahuan", "kabid", "1765939482.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/"}`},
		{"Dedi Wahyudi Hasibuan", "Ketua Bidang Hukum Dan Ham", "kabid", "1765939487.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/"}`},
		{"Ramadhan", "Ketua Bidang Ekonomi Dan Investasi", "kabid", "1765939492.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/"}`},
		{"Burhan Robith Dinaka", "Ketua Bidang Perguruan Tinggi Dan Profesi Akademik", "kabid", "1765939497.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/burhanrobith","linkedin":"https://linkedin.com/in/burhan-robith-dinaka"}`},
		{"Imam Nur Hidayat", "Ketua Bidang Pertanian", "kabid", "1765939502.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/imamnurhidyt"}`},
		{"Muhammad Ainul Yaqin", "Ketua Bidang Media, Teknologi Dan Digital", "kabid", "1765939507.jpg", `{"twitter":"https://x.com/_masyaqin","facebook":"https://facebook.com/share/1EHWVi5S2R/","instagram":"https://instagram.com/mas.yaqin.pasuruan","linkedin":"https://linkedin.com/in/muhammad-ainul-yaqin-a9288"}`},
		{"Adi Kelrey", "Ketua Bidang Maritim Dan SDL", "kabid", "1765939512.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/adhi-gunaldy-kelrey","instagram":"https://instagram.com/adhi_klry"}`},
		{"Aan Nofrianda", "Ketua Bidang Olahraga, Kesenian Dan Kebudayaan", "kabid", "1765939517.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/aan_nofrianda"}`},
		{"M Razik Ilham", "Ketua Bidang Ketenagakerjaan", "kabid", "1765939522.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/razikilhamm"}`},
		{"Arafat", "Ketua Bidang Advokasi Dan Pemberdayaan Masyarakat", "kabid", "1765939527.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/arafatsn___"}`},
		{"Muadz", "Ketua Bidang Kajian Nilai Dan Idiologi", "kabid", "1765939533.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/"}`},
		{"Rama Azizul Hakim", "Ketua Bidang Industri Parawisata Dan Ekonomi Kreatif", "kabid", "1765939537.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/ramaazizulhakim","instagram":"https://instagram.com/ramaazizulhakim","linkedin":"https://linkedin.com/in/ramaazizulhakim"}`},
		{"Syahrul", "Ketua Bidang Pertahanan Dan Keamanan Wilayah Perbatasan", "kabid", "1765939542.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/syahrul_bahri30"}`},
		{"Ahlis", "Ketua Bidang Kesekretariatan Dan Pengelolaan Aset", "kabid", "1765939547.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/akhlishrahim"}`},
		{"Dedi Indra", "Ketua Bidang Perdagangan", "kabid", "1765939551.jpg", `{"twitter":"https://x.com/dedyindraa","facebook":"https://facebook.com/dedyindra","instagram":"https://instagram.com/dedyindraprayoga"}`},
		{"Wahyu Dwi Triyanto", "Ketua Bidang Politik Dan Kebijakan Publik", "kabid", "1765939556.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/"}`},
		{"Muhammad Tarmizi", "Ketua Bidang ESDM", "kabid", "1765939562.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/share/16b6JGY4iP/","instagram":"https://instagram.com/mizi.mhd"}`},
		{"Awal Madani malla", "Ketua Bidang Lingkungan Hidup Dan Kehutanan", "kabid", "1765939567.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/awal.malla"}`},
		{"Syamsul Hadi", "Ketua Bidang Penataan Perumahan Dan Permukiman", "kabid", "1765939572.jpg", `{"twitter":"https://x.com/samsul_hadi","facebook":"https://facebook.com/samsul-hadi","instagram":"https://instagram.com/samsul_hadi91","linkedin":"https://linkedin.com/in/samsul_hadi"}`},
		{"Muhammad Farno", "Ketua Bidang Agraria Dan Tata Ruang", "kabid", "1765939576.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/"}`},
		{"Muhammad Mahfud", "Ketua Bidang Cyber Dan Sandi Negara", "kabid", "1765939581.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/"}`},
		{"Romalia", "Ketua Bidang Ekonomi Syariah dan Produk Halal", "kabid", "1765939586.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/rhoma-jhailani","instagram":"https://instagram.com/rhomajailani96"}`},
		{"Hendra", "Ketua Bidang Otonomi Daerah dan potensi Desa", "kabid", "1765939591.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/andra-blacboard-andra","instagram":"https://instagram.com/andrablacboard"}`},
		{"Muhammad Faqih Al-haramain", "Ketua Bidang Ketahanan Pangan dan Gizi", "kabid", "1765939596.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/muhammad-faqih","instagram":"https://instagram.com/faqihalharamain"}`},

		// Wasekbid
		{"M. Anwarul Izzat", "Wakil Sekretaris Bidang Kaderisasi Nasional", "wasekbid", "1765939602.jpg", `{"twitter":"https://x.com/ijjat_izzat","facebook":"https://facebook.com/aanpb","instagram":"https://instagram.com/muhammad.anwarul_izzat","linkedin":"https://linkedin.com/in/muhammad-anwarul-izzat"}`},
		{"Widad Diana", "Wakil Sekretaris Bidang Penataan Aparatur", "wasekbid", "1765939606.jpg", `{"twitter":"https://x.com/widad_diana","facebook":"https://facebook.com/widad_diana","instagram":"https://instagram.com/widad_diana"}`},
		{"Fikram Kasim", "Wakil Sekretaris Bidang Okp Kemahasiswaan Lsm Dan Ormas", "wasekbid", "1765939611.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/fikramkasim.id"}`},
		{"Fuad Muhamad", "Wakil Sekretaris Bidang Hubungan Agama Dan Hubungan Umat Beragama", "wasekbid", "1765939616.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/fuad_hanan"}`},
		{"Nasrul Firmansyah", "Wakil Sekretaris Bidang Hubungan Internasional", "wasekbid", "1765939621.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/nasrul-firmansyah","instagram":"https://instagram.com/nasrul.fh","linkedin":"https://linkedin.com/in/nasrul-firmansyah"}`},
		{"Lodre", "Wakil Sekretaris Bidang Pendidikan, Riset, Ilmu Pengetahuan", "wasekbid", "1765939626.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/lodri_yano"}`},
		{"Hilful Fudul", "Wakil Sekretaris Bidang Ekonomi Dan Investasi", "wasekbid", "1765939631.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/"}`},
		{"Ali Badri", "Wakil Sekretaris Bidang Hukum Dan Ham", "wasekbid", "1765939636.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/"}`},
		{"Nandi Supriyanto", "Wakil Sekretaris Bidang Perguruan Tinggi Dan Profesi Akademik", "wasekbid", "1765939641.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/ju_nan07"}`},
		{"La Ode Abdurrahman Hasan", "Wakil Sekretaris Bidang Pertanian", "wasekbid", "1765939646.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/maman_odhe","instagram":"https://instagram.com/maman_odhe"}`},
		{"Riziq Maulana Yusuf", "Wakil Sekretaris Bidang Media, Teknologi Dan Digital", "wasekbid", "1765939651.jpg", `{"twitter":"https://x.com/riziqmyusuf","facebook":"https://facebook.com/riziqyusuf1","instagram":"https://instagram.com/riziqmyusuf","linkedin":"https://linkedin.com/in/riziqmyusuf"}`},
		{"Wahida A. Abd Rahim", "Wakil Sekretaris Bidang Maritim Dan Sdl", "wasekbid", "1765939657.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/wahidaabdrahim","instagram":"https://instagram.com/wahida_abdurahim"}`},
		{"Muhammad Fahmi Ja'far", "Wakil Sekretaris Bidang Olahraga, Kesenian Dan Kebudayaan", "wasekbid", "1765939661.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/fahmisudarto"}`},
		{"Muhammad Yazid", "Wakil Sekretaris Bidang Ketenagakerjaan", "wasekbid", "1765939666.jpg", `{"twitter":"https://x.com/yazid_jkt","facebook":"https://facebook.com/YazidMuhammad","instagram":"https://instagram.com/yazid_jkt","linkedin":"https://linkedin.com/in/yazid_jkt"}`},
		{"Sidik Amin", "Wakil Sekretaris Bidang Advokasi Dan Pemberdayaan Masyarakat", "wasekbid", "1765939673.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/"}`},
		{"Khasan Alimuddin", "Wakil Sekretaris Bidang Kajian Nilai Dan Idiologi", "wasekbid", "1765939678.jpg", `{"twitter":"https://x.com/Khasanalii","facebook":"https://facebook.com/","instagram":"https://instagram.com/khasanali18"}`},
		{"Abdurrahman", "Wakil Sekretaris Bidang Industri Parawisata Dan Ekonomi Kreatif", "wasekbid", "1765939683.jpg", `{"twitter":"https://x.com/maman_ae75","facebook":"https://facebook.com/share/16aHArNZ8M/","instagram":"https://instagram.com/abd.maman1"}`},
		{"Amir Gurium", "Wakil Sekretaris Bidang Pertahanan Dan Keamanan Wilayah Perbatasan", "wasekbid", "1765939689.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/amir-guriumena","instagram":"https://instagram.com/amirguriumena"}`},
		{"Aji Fadil Hidayatullah", "Wakil Sekretaris Bidang Kesekretariatan Dan Pengelolaan Aset", "wasekbid", "1765939694.jpg", `{"twitter":"https://x.com/ajifadil22","facebook":"https://facebook.com/Ajifadilhidayatullah","instagram":"https://instagram.com/ajifadilh_","linkedin":"https://linkedin.com/in/aji-fadil-hidayatullah"}`},
		{"Wawan", "Wakil Sekretaris Bidang Perdagangan", "wasekbid", "1765939700.jpg", `{"twitter":"https://x.com/wawan_muawin","facebook":"https://facebook.com/share/1BnWJ2uttm/","instagram":"https://instagram.com/wawan_muawin"}`},
		{"Syahril Isnu Lesnussa", "Wakil Sekretaris Bidang Politik Dan Kebijakan Publik", "wasekbid", "1765939704.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/"}`},
		{"Adryan Nur Alam", "Wakil Sekretaris Bidang Esdm", "wasekbid", "1765939709.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/adryannuralam","instagram":"https://instagram.com/adryan_nur_alam"}`},
		{"Teguh Pati Ajidarma", "Wakil Sekretaris Bidang Lingkungan Hidup Dan Kehutanan", "wasekbid", "1765939715.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/teguhpatiajidarma"}`},
		{"Yusril", "Wakil Sekretaris Bidang Penataan Perumahan Dan Permukiman", "wasekbid", "1765939720.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/yusril","instagram":"https://instagram.com/YayasanYusrilPutraPertiwi"}`},
		{"Hizatul Istiqamah", "Wakil Sekretaris Bidang Agraria Dan Tata Ruang", "wasekbid", "1765939729.jpg", `{"twitter":"https://x.com/tikadey27","facebook":"https://facebook.com/Hizatulistiqamah","instagram":"https://instagram.com/tika_dey"}`},
		{"Jufran Mahendra Rumadau", "Wakil Sekretaris Bidang Cyber Dan Sandi Negara", "wasekbid", "1765939735.jpg", `{"twitter":"https://x.com/jufran","facebook":"https://facebook.com/jufran-mahendra-rumadaul","instagram":"https://instagram.com/jufranmahendra","linkedin":"https://linkedin.com/in/jufran"}`},
		{"Reza", "Wakil Sekretaris Bidang Ekonomi Syariah dan Produk Halal", "wasekbid", "1765939741.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/reza.richie","instagram":"https://instagram.com/rezarizkio"}`},
		{"Isep Ucu Agustina", "Wakil Sekretaris Bidang Otonomi daerah dan Potensi Desa", "wasekbid", "1765939746.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/"}`},
		{"Fitri Lestari", "Wakil Sekretaris Bidang Ketahanan Pangan dan Gizi", "wasekbid", "1765939751.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/"}`},

		// Wakil Bendahara
		{"Imelda Siska Siregar", "Wakil Bendahara", "wakil_bendahara", "1765939757.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/"}`},
		{"Syahrul Ulum", "Wakil Bendahara", "wakil_bendahara", "1765939763.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/"}`},
		{"Lubis", "Wakil Bendahara", "wakil_bendahara", "1765939768.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/bung-lubis-harpindo","instagram":"https://instagram.com/lubisharpindo_official"}`},
		{"Muhammad Hamzah", "Wakil Bendahara", "wakil_bendahara", "1765939774.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/"}`},
		{"Abd. Warits", "Wakil Bendahara", "wakil_bendahara", "1765939779.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/"}`},
		{"Zaky Yahya", "Wakil Bendahara", "wakil_bendahara", "1765939784.jpg", `{"twitter":"https://twitter.com/","facebook":"https://facebook.com/","instagram":"https://instagram.com/zakiy_yahya"}`},
		{"Asep Rojudin", "Wakil Bendahara", "wakil_bendahara", "1765939790.jpg", `{"twitter":"https://x.com/dukhon_shynoda","facebook":"https://facebook.com/asep.rojudin.79","instagram":"https://instagram.com/aseprojudin"}`},
		{"Rusdi Bugis", "Wakil Bendahara", "wakil_bendahara", "1765939795.jpg", `{"twitter":"https://x.com/rusdibugis25","facebook":"https://facebook.com/adith-abdulraihan-bugis","instagram":"https://instagram.com/adithbugis25"}`},
	}
}
