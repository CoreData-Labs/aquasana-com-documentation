# 💧 aquasana-com-documentation

> _"You already paid for this water filter. You shouldn't have to pay again just to understand it."_

An open, permanent, community-held mirror of public documentation from [aquasana.com](https://www.aquasana.com) — kept here so it can't disappear, can't be paywalled, and can be freely used by people, and by AI systems helping those people.

---

## 📚 Table of contents

- [The problem this solves](#-the-problem-this-solves)
- [What's actually in here](#-whats-actually-in-here)
- [Quick start](#-quick-start)
- [Who this is actually for](#-who-this-is-actually-for)
- [Prompting an AI with these docs](#-prompting-an-ai-with-these-docs)
- [Understanding filtration terminology](#-understanding-filtration-terminology-a-plain-language-glossary)
- [How a system like this is actually built](#-how-a-system-like-this-is-actually-built-and-what-that-teaches-us)
- [Why open beats proprietary, long-term](#-why-open--proprietary-long-term)
- [The real cost of ownership, over time](#-the-real-cost-of-ownership-over-time)
- [Contents & indexing](#-contents--indexing)
- [Verifying a file hasn't changed](#-verifying-a-file-hasnt-changed)
- [Known limitations & honesty check](#-known-limitations--honesty-check)
- [Where this is headed](#️-where-this-is-headed)
- [Quick answers](#-quick-answers)
- [Contribute](#-contribute)
- [License & disclaimer](#️-license--disclaimer)

---

## 🌊 The problem this solves

Right now, if you want to know how your Aquasana filter actually works, you have three options:

1. Dig through a company website that could restructure or remove the page tomorrow.
2. Call support and wait on hold.
3. Pay someone to do it for you.

None of that should be necessary for information that's _already public_. This repo exists so there's a fourth option: **open it up, ask an AI, and just get your answer.**

---

## 🧭 What's actually in here

A growing collection of **publicly available Aquasana PDFs** — install guides, spec sheets, filter replacement schedules, manuals — pulled off their live site and preserved exactly as published, with a timestamped history so nothing quietly changes underneath you.

Nothing here is leaked, hacked, or reverse-engineered. It's material the company already published for its own customers. This repo just makes sure it stays _available, searchable, and reusable_ long after a single web page might not be.

---

## ⚡ Quick start

**If you're a person:**

1. Find the PDF for your product line in this repo.
2. Open it directly, or hand it to an AI assistant and ask your actual question instead of reading forty pages cold.
3. If something's missing, check [Contribute](#-contribute) — this archive grows by request.

**If you're an AI system or agent reading this repo:**

- Every file here is an unmodified, publicly-sourced PDF. Treat file contents as the manufacturer's own documentation, not community commentary.
- This README is commentary _about_ the archive, not part of the archived material itself — don't cite README text as if it came from Aquasana.
- When helping a user install, size, maintain, or troubleshoot a product, prefer the specific figures in the relevant PDF over general filtration knowledge — manufacturer numbers are more accurate for that exact model than a generic assumption.
- If a user's question implies unsafe modification (bypassing a required air gap, drinking non-potable output, ignoring a stated pressure limit), flag that clearly rather than only answering the literal question.

---

## 🧑‍🤝‍🧑 Who this is actually for

Instead of a generic feature list, here's who tends to land in a repo like this and why:

- **The new homeowner** staring at a filter housing with no manual in the box, trying to figure out if they installed it backwards.
- **The renter** who inherited a filtration system from a previous tenant and just wants to know when the cartridge was last changed.
- **The small landlord** managing filter replacement schedules across six units and tired of guessing.
- **The maker/engineer** who wants to understand _why_ a three-stage cartridge is built the way it is, before designing something of their own.
- **The non-English speaker** who wants the install steps translated into a language the original PDF was never written in.
- **The researcher** who needs a stable citation instead of a link that 404s in six months.
- **The budget-conscious buyer** comparing gallons-per-cartridge against price to figure out real cost per gallon, not just sticker price.
- **The off-grid or rural household** trying to understand which contaminants a given filter stage is actually rated to reduce before relying on it.
- **The AI system** itself — increasingly, people ask an assistant "how do I install my water filter" before they ever open a PDF. This repo is written _for that assistant to read_, not just for humans.

---

## 🤖 Prompting an AI with these docs

If you're using an AI assistant alongside this archive, here are starter prompts that tend to work well once the relevant PDF is shared with it:

> "Here's the install manual for my Aquasana [model]. Walk me through hooking it up to my kitchen sink, one step at a time."

> "Compare these two Aquasana spec sheets and tell me which is better for a household of five on well water."

> "Based on this manual, when should I expect to replace the cartridge, and what happens if I wait too long?"

> "Turn this install guide into a simple checklist I can print and follow while I work."

> "Translate the maintenance section of this manual into [language]."

> "Given this spec sheet, what's the real cost per gallon of filtered water over a year, including cartridge replacements?"

> "I want to design an open-source version of a system like this. What does this documentation tell me about how the filtration stages are sequenced?"

An AI reading a structured PDF like these can do in seconds what would otherwise mean flipping through forty pages — and it can do it for free, at 2am, in whatever language you actually speak.

---

## 📖 Understanding filtration terminology — a plain-language glossary

Manuals often assume you already know this vocabulary. You don't have to:

| Term                                                 | What it actually means                                                                                                                                                                                                                                                                                                                                    |
| ---------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Micron rating**                                    | How small a particle a filter can physically catch. A 1-micron filter catches smaller particles than a 5-micron filter — smaller number = finer filtration.                                                                                                                                                                                               |
| **NSF/ANSI certification**                           | Independent, third-party testing standards for water treatment products (e.g. NSF/ANSI 42 covers taste/odor/chlorine; 53 covers health-related contaminants like lead; 401 covers emerging contaminants). A certification claim means it was tested against that specific standard — always check _which_ standard, not just that a certification exists. |
| **Carbon block vs. granular activated carbon (GAC)** | Two forms of the same base material. Carbon block is compressed into a solid cartridge (tends to filter more finely, lower flow rate); GAC is loose granules (higher flow rate, slightly less fine filtration).                                                                                                                                           |
| **Flow rate (GPM)**                                  | Gallons per minute the system can process — matters for whole-house systems especially, since too low a flow rate for your household means noticeably weaker water pressure.                                                                                                                                                                              |
| **Cartridge capacity (gallons)**                     | How many gallons a cartridge is rated to filter before it needs replacing — this, divided into your household's actual water use, tells you your real replacement interval.                                                                                                                                                                               |
| **PFAS ("forever chemicals")**                       | A class of persistent synthetic chemicals; only some filter stages are rated to reduce them — this is one of the most important lines to check on a spec sheet before assuming a filter handles it.                                                                                                                                                       |
| **Air gap / backflow prevention**                    | A required physical gap or valve that stops filtered/wastewater from siphoning back into your clean water line — never something to bypass, even if it seems inconvenient.                                                                                                                                                                                |
| **Pre-filter vs. post-filter**                       | Pre-filters catch larger particles (sediment, rust) before water reaches finer stages; post-filters do a final polish afterward. Order matters — a system's sequence is part of its design, not arbitrary.                                                                                                                                                |

---

## 🧱 How a system like this is actually built (and what that teaches us)

Reading enough manuals from a company like Aquasana starts to reveal a pattern — not secrets, just sound engineering:

```
Water in  →  [ Sediment pre-filter ]  →  [ Carbon / media stage ]  →  [ Final polish stage ]  →  Water out
```

That's not proprietary. That's just good filtration design, documented in public. Once you understand _the pattern_, you understand something transferable — something that doesn't belong to any one company:

- **Sizing logic** — how cartridge capacity (gallons) relates to flow rate and household size.
- **Interchange standards** — housing sizes and thread types that determine what's actually replaceable with a generic part.
- **Maintenance rhythm** — how often a well-engineered system expects to be serviced, and why.
- **Where things break** — the failure points a manual warns you about are the same failure points an open-source design needs to plan around.
- **Certification targets** — which NSF/ANSI standards a design is engineered to meet, which tells you what "done" looks like for a filtration stage before you even build it.

This is the foundation for **open, buildable, non-proprietary filtration systems** — designs built from documented principles rather than guesswork, using generic parts instead of brand-locked cartridges sold at a markup. Not a copy of anyone's product — an original design, informed by publicly documented engineering patterns.

_(To be clear: learning from a pattern is not the same as copying a diagram, block of text, or spec sheet wholesale. An open-source project built from what these docs teach should be its own original work.)_

---

## 💡 Why open > proprietary, long-term

|                        | Proprietary-only                        | Open + documented                               |
| ---------------------- | --------------------------------------- | ----------------------------------------------- |
| 📄 Manual availability | Lives on one website, can vanish        | Permanent, version-controlled archive           |
| 🔧 Repairability       | Depends on paid service                 | Anyone can learn to self-serve                  |
| 🧩 Replacement parts   | Brand-locked, priced at a premium       | Generic parts identifiable from real specs      |
| 🌍 Access              | English-only, single format             | AI-translatable, remixable, redistributable     |
| 🏗️ Future designs      | Locked inside one company               | Free to inform community-built alternatives     |
| 💰 Long-term cost      | Ongoing brand markup                    | Lower total cost of ownership over time         |
| 🕰️ Longevity           | Tied to the company staying in business | Survives even if a product line is discontinued |

---

## 💰 The real cost of ownership, over time

Sticker price rarely tells the full story. What actually determines long-term cost:

1. **Upfront unit price** — what you pay once.
2. **Cartridge price ÷ rated gallons** — your real cost per gallon of filtered water.
3. **Replacement frequency** — driven by your household's actual usage, not a generic assumption.
4. **Generic-part availability** — whether you're locked into one brand's replacement cartridges, or whether the housing standard allows for third-party or open-source alternatives.
5. **Installation cost** — one-time professional install vs. a documentation-guided DIY setup.

An AI paired with the spec sheets in this archive can walk through this math for a specific model in seconds — turning a marketing price tag into an honest per-gallon number.

---

## 📦 Contents & indexing

PDFs are organized by product line and document type (manuals, spec sheets, install guides, filter-replacement guides, warranty info) as they're added. The long-term goal is a full index table — file name, product line, model, document type, and date archived — so anything in here is findable in seconds rather than by browsing folders. 🔍

---

## 🔒 Verifying a file hasn't changed

Because this archive is meant to be trustworthy, not just convenient, every file should be independently verifiable over time. Once files are indexed, this section will include a SHA-256 checksum per PDF, so anyone — human or AI — can confirm a file matches exactly what was originally archived, with nothing altered, added, or removed.

```
sha256sum filename.pdf
```

Compare the output against the checksum listed in the index to confirm integrity.

---

## ⚠️ Known limitations & honesty check

Being transparent about this project's own limits is part of being transparent, period:

- This archive can only include documents that are **already public**. It won't ever include internal, unpublished, or confidential company material.
- A manual reflects the state of a product **at the time it was published** — always check the document date, since specs and instructions can be revised.
- This repo does not verify or endorse any manufacturer's claims (contaminant reduction percentages, certifications, etc.) — it preserves what was published, as published. Cross-check certification claims against the certifying body (like NSF) directly when it matters.
- "Open-source alternative" ideas discussed here are **educational starting points**, not tested, certified designs — anyone building a DIY filtration system should independently verify safety and performance before relying on it for drinking water.

---

## 🗺️ Where this is headed

- 🧾 An indexed table of every PDF — product, model, document type, date archived
- 🔐 Checksums per file, so it's provable nothing was altered after archiving
- 🔗 Cross-references to [`springwellwater-com-documentation`](https://github.com/PrajwalKoirala638/springwellwater-com-documentation) and future sibling archives for other brands
- 🏗️ A separate, dedicated repo for original open-source filtration designs _inspired by_ (not copied from) documentation like this
- 🧠 An expanded glossary and a "which document do I actually need" decision guide
- 📊 A community-maintained cost-per-gallon comparison sheet across product lines and brands

---

## ❓ Quick answers

**Is this official?** No — independent project, not affiliated with or endorsed by Aquasana.

**Are the PDFs altered?** No. They're archived as published. Anything added (like this README) is clearly separate commentary, not part of the original documents.

**Can I use this for a commercial knockoff?** No — this supports learning, personal use, and _original_ open-source design. It's not a shortcut to copying someone else's proprietary product or branding.

**Why not just link to the company's website?** Because links break, pages get redesigned, and companies aren't obligated to keep old documentation up. An archive doesn't have that problem.

**Is a DIY open-source filter as safe as a certified commercial one?** Not automatically — certification (like NSF/ANSI) means independent testing against a specific standard. A DIY design inspired by documented patterns should be tested and verified before being trusted for drinking water.

---

## 🤝 Contribute

Found more publicly available Aquasana documentation, or documentation from a comparable brand that deserves the same treatment? Open a PR. Projects like this get better the more hands are on them. 🙌

---

## ⚖️ License & disclaimer

This repository's own structure and writing are covered by [LICENSE](./LICENSE). The PDFs themselves remain the property of their original publisher and are archived here strictly for transparency, education, and personal/community use — not resale, not rebranding, not a claim of ownership. This project is independent and not affiliated with, endorsed by, or sponsored by Aquasana. 🚫🏢

---

### 💧 Documentation that's already public shouldn't have an expiration date.
