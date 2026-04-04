import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Политика за приватност | плата.fyi",
};

export default function PrivacyPage() {
  return (
    <div className="max-w-2xl mx-auto py-12 space-y-16">
      <div>
        <h1 className="text-4xl sm:text-6xl font-black tracking-tight leading-none text-ink">
          Како работи приватноста?
        </h1>
        <p className="mt-4 text-lg font-medium opacity-60 max-w-lg">
          Без е-маил, без сметки, без лозинки.
        </p>
      </div>

      <section className="space-y-3">
        <h2 className="text-2xl font-black text-ink">Кои податоци се собираат?</h2>
        <p className="text-base opacity-70 leading-relaxed">
          Plata.fyi собира само доброволно споделени податоци за плата:
          компанија, позиција, индустрија, град, искуство и примања. Сите
          јавно прикажани записи се целосно анонимни.
        </p>
      </section>

      <section id="anonimnost" className="space-y-3" style={{ scrollMarginTop: "72px" }}>
        <h2 className="text-2xl font-black text-ink">Kako функционира анонимноста?</h2>
        <p className="text-base opacity-70 leading-relaxed">
          Без е-маил, без лозинка, без регистрација. Пополнуваш форма, поминуваш кратка верификација
          (Cloudflare Turnstile) и добиваш анонимна сесија зачувана во твојот прелистувач. Важи 30 дена.
        </p>
        <p className="text-base opacity-70 leading-relaxed">
          Ако го смениш уредот, прелистувачот, или ги избришеш податоците, старите записи нема да бидат
          достапни — бидејќи не постои е-маил или акаунт преку кој би се вратил до нив. Ова е намерно.
        </p>
        <h4 className="text-l font-black text-ink">Kako анонимноста функционира технички?</h4>
        <ol className="space-y-2 text-base opacity-70 leading-relaxed list-decimal list-inside">
          <li>Пополнуваш форма за плата.</li>
          <li>Поминуваш верификација (Cloudflare Turnstile) — докажува дека си човек, не собира лични податоци.</li>
          <li>Серверот генерира анонимен токен и го зачувува во твојот прелистувач (localStorage). Важи 30 дена.</li>
          <li>Додека сесијата е активна, можеш да ги уредуваш и бришеш твоите записи без повторно пријавување.</li>
          <li>Ако изгубиш пристап (нов уред, нов прелистувач, избришани податоци), можеш да поднесеш нов запис — но тој ќе биде под нова, празна сесија. <em>Старите записи не можат да се вратат.</em></li>
        </ol>
        <p className="text-base opacity-70 leading-relaxed">
          Никаква е-маил адреса не се собира или зачувува на никаков начин. Не постои врска помеѓу тебе
          и твоите записи — дури ни ние не можеме да утврдиме кој ги внел. Датумот на поднесување се
          зачувува само со прецизност на месец (пр. Април 2026), без точен час или ден.
        </p>
      </section>

      <section className="space-y-3">
        <h2 className="text-2xl font-black text-ink">Јавни податоци</h2>
        <p className="text-base opacity-70 leading-relaxed">
          Сите записи за плата се јавно достапни и анонимни. Немаме
          можност да утврдиме кој ги внел одредени податоци врз основа на
          јавно достапните информации.
        </p>
      </section>

      <section className="space-y-3">
        <h2 className="text-2xl font-black text-ink">Повратни информации и предлози</h2>
        <p className="text-base opacity-70 leading-relaxed">
          Имаш идеја за подобрување, забележа грешка, или сакаш да ни кажеш нешто? Пишете на{" "}
          <a href="mailto:kontakt@plata.fyi" className="underline underline-offset-2 hover:opacity-100 transition-opacity">
            kontakt@plata.fyi
          </a>
          .
        </p>
      </section>

      <section className="space-y-3 pb-4 border-t-2 border-ink/10 pt-10">
        <h2 className="text-2xl font-black text-ink">Инспирација</h2>
        <p className="text-base opacity-70 leading-relaxed">
          Plata.fyi е инспириран од{" "}
          <a href="https://levels.fyi" target="_blank" rel="noopener noreferrer" className="underline underline-offset-2 hover:opacity-100 transition-opacity">
            levels.fyi
          </a>
          , платформа која им овозможува на вработените во технолошката индустрија да ги споредуваат своите компензации. Целта е да се направи исто за сите индустрии во Македонија.
        </p>
      </section>
    </div>
  );
}