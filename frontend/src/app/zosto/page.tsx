import Link from "next/link";

export const metadata = {
  title: "Зошто да споделам? | плата.fyi",
  description:
    "Зошто да ја споделиш твојата плата? Транспарентност, статистика и поддршка при преговарање, без комерцијални интереси.",
};

function ShareButton() {
  return (
    <Link href="/submit" className="btn-primary inline-block mt-6">
      Споделете плата →
    </Link>
  );
}

export default function ZostoPage() {
  return (
    <div className="max-w-2xl mx-auto py-12 space-y-16">
      {/* Intro */}
      <div>
        <h1 className="text-4xl sm:text-6xl font-black tracking-tight leading-none text-ink">
          Твојата плата<br /> има вредност.
        </h1>
        <p className="mt-4 text-lg font-medium opacity-60 max-w-lg">
          <strong><i>Еден бројка е анегдота, илјадници бројки се статистика.</i></strong>
        </p>
      </div>

      {/* Section 1 */}
      <section className="space-y-3">
        <h2 className="text-2xl font-black text-ink">За статистиката</h2>
        <p className="text-base opacity-70 leading-relaxed">
          Во Македонија речиси нема јавни и веродостојни податоци за платите.
          Државните статистики се агрегирани и застарени. Само со масовно, доброволно споделување
          можеме да изградиме реална слика по индустрија, по град, по ниво на искуство.
        </p>
        <p className="text-base opacity-70 leading-relaxed">
          Секој запис што го додаваш го зголемува квалитетот на статистиката за сите.
          Колку повеќе луѓе учествуваат, толку попрецизни се просеците.
        </p>
        <ShareButton />
      </section>

      {/* Section 2 */}
      <section className="space-y-3">
        <h2 className="text-2xl font-black text-ink">За транспарентноста</h2>
        <p className="text-base opacity-70 leading-relaxed">
          Немањето информации е во корист на работодавачите. Тие знаат колку плаќаат
          на сите позиции, а ти знаеш само за својата.
        </p>

        <div className="space-y-2 pt-2">
          <h3 className="text-base font-bold text-ink opacity-80">Зошто ова е важно</h3>
          <table className="w-full text-sm border-collapse" style={{ border: "2px solid rgb(40,40,37)" }}>
            <thead>
              <tr style={{ borderBottom: "2px solid rgb(40,40,37)" }}>
                <th className="text-left px-4 py-2 font-black text-ink w-1/2" style={{ borderRight: "2px solid rgb(40,40,37)" }}>
                  Кога има помалку податоци
                </th>
                <th className="text-left px-4 py-2 font-black text-ink w-1/2">
                  Кога има повеќе податоци
                </th>
              </tr>
            </thead>
            <tbody className="opacity-70">
              <tr style={{ borderBottom: "1px solid rgba(40,40,37,0.15)" }}>
                <td className="px-4 py-2" style={{ borderRight: "2px solid rgb(40,40,37)" }}>полесно се создаваат нерамнотежи во плати</td>
                <td className="px-4 py-2">очекувањата се пореални</td>
              </tr>
              <tr style={{ borderBottom: "1px solid rgba(40,40,37,0.15)" }}>
                <td className="px-4 py-2" style={{ borderRight: "2px solid rgb(40,40,37)" }}>луѓето потешко преговараат</td>
                <td className="px-4 py-2">преговорите се пофер</td>
              </tr>
              <tr>
                <td className="px-4 py-2" style={{ borderRight: "2px solid rgb(40,40,37)" }}>одлуките се базираат на претпоставки</td>
                <td className="px-4 py-2">пазарот станува појасен за сите</td>
              </tr>
            </tbody>
          </table>
        </div>

        <p className="text-base opacity-70 leading-relaxed">
          Транспарентноста не те прави ранлив, те прави информиран. Анонимноста е
          вградена во веб-сајтот. <Link href="/privacy" className="underline underline-offset-2 hover:opacity-100 transition-opacity">Како работи анонимоста?</Link>
        </p>
      </section>

      {/* Section 4 */}
      <section className="space-y-3">
        <h2 className="text-2xl font-black text-ink">Алатка за преговарање</h2>
        <p className="text-base opacity-70 leading-relaxed">
          Кога одиш на интервју или барање за зголемување, добро е да дојдеш со факти.
          „Јас сметам дека заслужувам повеќе" е поинаква појдовна точка од „просекот за оваа
          работа е X, јас барам Y."
        </p>
        <p className="text-base opacity-70 leading-relaxed">
          Колку побогата базата на податоци, толку посилен е аргументот за секој поединечно.
        </p>
      </section>

      {/* Section 5 */}
      <section className="space-y-3 pb-4 border-t-2 border-ink/10 pt-10">
        <h2 className="text-2xl font-black text-ink">Без комерцијални интереси</h2>
        <p className="text-base opacity-70 leading-relaxed">
          плата.fyi не е бизнис. Нема реклами, нема претплата, нема продавање на
          податоци. Идејата е едноставна: луѓето да имаат референтна точка.
        </p>
        <p className="text-base opacity-70 leading-relaxed">
          Платформата постои како јавен, слободен ресурс за сите кои сакаат да знаат каде стојат.
        </p>
        <ShareButton />
      </section>
    </div>
  );
}
