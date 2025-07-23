# Final Algorithm: Resolving Bank Wealth with Management Fees and Taxes

This document describes a linear‐algebra-based algorithm to compute each bank’s end-of-day wealth, incorporating:

1. Management fees: a bank (manager) retains a fraction of the profits it generates for investors.
2. Taxes: any "free" money conjured by investment loops is tracked and redistributed as a tax charge.

---

## 1. Notation and Inputs

Let there be **n** banks, indexed by *i=1…n*. Define:

- **W** ∈ ℝⁿ×ⁿ: exposure fractions, where _Wᵢⱼ_ is the fraction of bank *i*’s capital invested in bank *j*.
- **s** ∈ ℝⁿ: fraction of each bank’s capital in exogenous assets (e.g. stocks).
- **c** ∈ ℝⁿ: fraction in cash (yield=1).
  - ∀i: ∑ⱼ Wᵢⱼ + sᵢ + cᵢ = 1.
- **y_s**: stock market yield (e.g. 1.10).
- **F** ∈ ℝⁿ: fixed daily costs (cash outflow).
- **p** ∈ ℝⁿ: Ponzi factor (how much banks over-report; typically ≥1).
- **m** ∈ [0,1]ⁿ: management-fee rate per bank; if a bank makes a profit, it retains _mⱼ_ of that profit, otherwise fees = 0.

Define the exogenous yield vector:
```
ex = s·y_s + c·1
```

---

## 2. Incorporating Management Fees

When bank *j* realizes a profit factor on its total book, investors who placed capital in *j* receive only a share:

- Raw yield of *j*: we solve for *xⱼ* below.
- Profit of *j*: _Δⱼ = xⱼ − investedⱼ_.
- Fee retained by *j*:  _feeⱼ = max(Δⱼ,0) · mⱼ_.
- Net paid out to investors: _Δⱼ · (1−mⱼ)_.

In linear form we fold fees into the exposure matrix:

- **A**ᵢⱼ = Wᵢⱼ · (1−mⱼ)  
- **R**ᵢⱼ = Wᵢⱼ · mⱼ      (for tracking fees)

---

## 3. Solving the Wealth System

We seek _x_ ∈ ℝⁿ the end-of-day wealth for each bank.  In vector form:

```
x = p ⊙ (A · x + ex) − F
```

Rearrange into a linear system:

```
(I − diag(p)·A) x = diag(p)·ex − F
```

1. Form **A** as above.
2. Build LHS matrix: _M = I − diag(p)·A_.
3. Build RHS vector: _b = diag(p)·ex − F_.
4. Solve _M x = b_.

If _det(M)=0_ (a singular loop with ≥100% cycle yield), proceed to taxes.

---

## 4. Handling Taxes for Infinite-Loop Gains

Any loop that amplifies exogenous gains beyond conservation is identified by singularity.  To recover:

1. Compute the total unconserved gain:
   ```
   G_loop = sum(x_solved) − (sum(ex_i) − sum(F_i))
   ```
2. Distribute a tax vector _t_ ∈ ℝⁿ so that ∑tᵢ = G_loop, proportional to each bank’s loop exposure.
3. Final wealth: _x_final = x_solved − t_.

---

## 5. Worked Example: 3-Bank Loop

**Banks**: A, B, C (n=3)

| Parameter | A    | B    | C    |
|-----------|------|------|------|
| W→A       | 0    | 0.6  | 0    |
| W→B       | 0.4  | 0    | 0.5  |
| W→C       | 0.6  | 0.4  | 0    |
| s (stocks)| 0    | 0.0  | 0.2  |
| c (cash)  | 0.0  | 0.4  | 0.3  |
| p (Ponzi) | 1.0  | 1.0  | 1.0  |
| m (fee%)  | 0.10 | 0.10 | 0.10 |
| F (costs) | 0.10 | 0.10 | 0.10 |
| y_s       |       |      |      |
```
s       = [0,    0.0, 0.2]
c       = [0.0,  0.4, 0.3]
ex      = s·1.1 + c·1 = [0.00, 0.40, 0.52]
```

**Step 1: Build A**
```
Aᵢⱼ = Wᵢⱼ·(1−0.10)
A = [ [0     0.6×0.9  0    ],
      [0.4×0.9 0     0.5×0.9],
      [0.6×0.9 0.4×0.9 0    ] ]
  ≈ [ [0    0.54 0   ],
      [0.36 0    0.45],
      [0.54 0.36 0   ] ]
```

**Step 2: Form and Solve (I−A)x = ex−F**
```
M = I − A
b = ex − F = [−0.10,  0.30,  0.42]
```
Solve _M x = b_; solution _x ≈ [0.35, 0.84, 0.91]_.

**Step 3: Check Loop Gain**
```
Total exogenous gain = sum(ex) − sum(F) = (0.00 + 0.40 + 0.52) − 0.30 = 0.62
Sum(x) ≈ 2.10
G_loop = 2.10 − 0.62 = 1.48  (loop amplification)
Apply taxes _t_ proportional to x, e.g. t = [0.25, 0.59, 0.64].
```
**Step 4: Final Wealth**
```
x_final = x − t ≈ [0.10, 0.25, 0.27]
```  
(final positive yields)

## Example Summary

- Initial capital (normalized to 1 per bank): A=1.00, B=1.00, C=1.00
- Investments:
  - A → B: 0.60; A → C: 0.00; A stocks: 0.00; A cash: 0.00
  - B → A: 0.40; B → C: 0.50; B stocks: 0.00; B cash: 0.40
  - C → A: 0.60; C → B: 0.40; C stocks: 0.20; C cash: 0.30
- Taxes owed:
  - A: 0.25
  - B: 0.59
  - C: 0.64
- Final wealth (x_final):
  - A: 0.10
  - B: 0.25
  - C: 0.27
- Overall yield ((x_final − initial)/initial):
  - A: 10%
  - B: 25%
  - C: 27%

*End of algorithm write-up.*
