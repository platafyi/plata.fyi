export default function Footer() {
  return (
    <footer className="bg-white mt-auto" style={{ borderTop: "2px solid rgb(40,40,37)" }}>
      <div className="container mx-auto px-4 max-w-6xl py-5">
        <div className="flex flex-col sm:flex-row justify-between items-center gap-3 text-sm">
          <p className="font-semibold text-ink opacity-60">
            © {new Date().getFullYear()} Plata.fyi, Сите податоци се анонимни
          </p>
          <div className="flex gap-4 font-semibold text-ink opacity-60">
            <a href="/privacy" className="hover:opacity-100 transition-opacity">Приватност</a>
            <a href="mailto:kontakt@plata.fyi" className="hover:opacity-100 transition-opacity">Контакт</a>
          </div>
        </div>
      </div>
    </footer>
  );
}
