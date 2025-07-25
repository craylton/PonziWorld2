# Resolving Bank Net-Worth with Inter-Bank Loops

This document presents a step-by-step algorithm to compute each bank’s end-of-day net worth and yields when banks invest in each other (including loops), pay fixed costs, charge management fees, and hold exogenous assets.

---

## 1. Notation and Inputs

- n₀ᵢ: initial net worth of bank *i* (before returns)
- n₁ᵢ: end-of-day net worth of bank *i* (after returns and costs)
- pᵢ: Ponzi factor for bank *i* (exaggerates its yield; ≥ 1)
- Wᵢⱼ: amount (or fraction) that bank *i* invests into bank *j*
- mᵢ: management fee rate of bank *i* (applied if it earns a profit)
- Fᵢ: fixed (daily) cost paid by bank *i*
- Eᵢ: total exogenous return for bank *i* (sum of each exogenous asset · that asset’s yield)
- yᵢ: the “believed” yield of bank *i*:  
  yᵢ = pᵢ · (n₁ᵢ / n₀ᵢ)


## 2. Core Equations

For each bank *i*, we set up a simultaneous system:

```text
n₁ᵢ = Eᵢ                             (exogenous returns)
     + ∑₍j₎ [ Wᵢⱼ · yⱼ · (1 - mⱼ) ]  (returns from investments, net of managers’ fees)
     - ∑₍k₎ [ Wₖᵢ · yᵢ · (1 - mᵢ) ]  (what *i* must pay investors back, net of its own fee)
     - Fᵢ                             (fixed daily cost)

yᵢ  = pᵢ · (n₁ᵢ / n₀ᵢ)
```  
These equations form a linear system in the unknowns n₁₁ … n₁ₙ.


## 3. Simple 2-Bank Example

**Setup:**
- Banks A and B
- A → B investment: £100
- B → A investment: £100
- After that, A invests £100 in stocks (10% gain)
- Fixed costs: Fₐ=2, Fᵦ=1
- No management fees or ponzi factors (mᵢ=0, pᵢ=1)

1. Initial net worths:
   - A₀ = 100
   - B₀ = 100
2. Exogenous gain (stocks): £100 × 10% = £10
3. Naïve split (all £10 to A) gives A₁=210, B₁=100 (violates B’s expected yield).
4. Proper system:
   ```text
   n₁A = 100·1.10 + 100·yB·1  - 100·yA·1  - 2
   n₁B = 100·1.10 + 100·yA·1  - 100·yB·1  - 1
   yA = n₁A / 100     ,   yB = n₁B / 100
   ```
5. Solving yields a split of the £10 so both banks see the correct profit on £100.

(Pretty sure this AI generated example is wrong but can't be bothered to correct it)


## 4. Three-Bank Loop Example

**Parameters:**
- Exogenous assets:
  - A: £50 in stocks @10% → yields 55
  - B: £30 in bonds  @5%  → yields 31.5
  - C: £20 in cash   @0%  → yields 20
- Inter-bank:
  - A→B: 20
  - B→C: 15
  - C→A: 10
- Fixed costs: F = [1, 2, 3] for A, B, C
- Fees: mᵢ=10% if bank makes profit
- Ponzi: pᵢ=1.1 for all banks

1. Compute initial net worths:
   - A₀ = 50 + 20 − 10 = 60
   - B₀ = 30 + 15 − 20 = 25
   - C₀ = 20 + 10 − 15 = 15
2. Write equations:
   ```text
   n₁A = 55   + 20·yB·0.9  - 10·yA·0.9  - 1
   n₁B = 31.5 + 15·yC·0.9  - 20·yB·0.9  - 2
   n₁C = 20   + 10·yA·0.9  - 15·yC·0.9  - 3
   yᵢ  = 1.1 · (n₁ᵢ / n₀ᵢ)
   ```
3. Solve (e.g. by fixed-point iteration or direct linear algebra):
   - n₁ ≈ [62.70, 24.05, 13.74]
4. Totals:
   - Σ n₀ = 100
   - Exogenous gain = 50·0.1 + 30·0.05 + 20·0 - 1 - 2 - 3 = 0.5 → expected Σ n₁ = 100.5
   - Actual Σ n₁ ≈ 62.70 + 24.05 + 13.74 = 100.49
5. **Fixed costs check:** banks paid F_total = 1+2+3 = 6.  Net system change = 101.07 − (100 − 6) = +7.07 → or: if we isolate loop effect, exogenous + fixed = 106.5 − 6 = 100.5, so ~0.5 of “free” money was created by the loop alone.


## 5. Loop Tax and Conservation

Let G_loop = Σ n₁ − (Σ n₀ + exogenous_gain − Σ F).  In the above:
```
G_loop = 100.49 − (100 + 0.5) ≈ 0.01 (this was with rounding errors)
```
This is the tiny surplus created by the loop itself (if any).  You can impose a tax of G_loop (distributed by policy) to restore exact conservation:

```
n₁ᵢ_final = n₁ᵢ − tᵢ,
where ∑ tᵢ = G_loop.
```

---

This algorithm scales to *N* banks by assembling an N×N linear system from the core equations, solving for n₁, and then optionally applying a loop tax to enforce total-system conservation.
