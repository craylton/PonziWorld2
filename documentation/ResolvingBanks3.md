# Resolving Bank Wealth with Management Fees and Taxes

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
- **m** ∈ [0,1]: management-fee rate applied to all banks; if a bank makes a profit, it retains _m_ of that profit, otherwise fees = 0.
- **k** ∈ ℝⁿ: initial capital per bank. Variables W, s, c are fractions of each bank's wealth. The solution x gives relative end-of-day wealth; actual wealth = k ⊙ x.

Define the exogenous yield vector:
```
ex = s·y_s + c·1
```

---

## 2. Incorporating Management Fees

When bank *j* realizes a profit factor on its total book, investors who placed capital in *j* receive only a share:

- Raw yield of *j*: we solve for *xⱼ* below.
- Profit of *j*: _Δⱼ = xⱼ − investedⱼ_.
- Fee retained by *j*:  _feeⱼ = max(Δⱼ,0) · m_.
- Net paid out to investors: _Δⱼ · (1−m)_.

In linear form we fold fees into the exposure matrix:

- **A**ᵢⱼ = Wᵢⱼ · (1−m)  
- **R**ᵢⱼ = Wᵢⱼ · m      (for tracking fees)

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
## 5. Simple Two-Bank Example

**Banks**: A, B (n=2)

| Parameter        | A     | B     |
|------------------|-------|-------|
| W→A              | 0     | 0     |
| W→B              | 0.5   | 0     |
| s (stocks)       | 0     | 0.5   |
| c (cash)         | 0.5   | 0.5   |
| p (Ponzi)        | 1.1   | 1.0   |
| F (costs)        | 0.05  | 0.10  |
| m (fee%)         | 0.20  | 0.20  |
| y_s (stock yield)| 1.20  | 1.20  |

### 5.1. Initialization and Inputs

We normalize initial wealth to 1 per bank, then scale later.  Define:
```
W = [ [0,   0.5],   # A invests 50% of its capital in B
      [0,   0  ] ]  # B invests nothing in A or itself
s = [0, 0.5]
c = [0.5, 0.5]
p = [1.1, 1.0]
m = 0.20
F = [0.05, 0.10]
y_s = 1.20
```
Exogenous yields for each bank:
```
ex = s·y_s + c·1 = [0*1.20 + 0.5*1, 0.5*1.20 + 0.5*1]
   = [0.50, 1.10]
```

### 5.2. Fee-Adjusted Exposure Matrix

Fees reduce the fraction paid out.  Compute A = W * (1−m):
```
A₁₂ = 0.5 * (1−0.20) = 0.5 * 0.8 = 0.40
All other Aᵢⱼ = 0

A = [ [0, 0.40],
      [0, 0   ] ]
```

### 5.3. Forming the Linear System

We solve (I − diag(p)·A)·x = diag(p)·ex − F.  Write M and b:
```
M = I − diag(p)·A = [ [1 − 1.1·0,   −1.1·0.40],
                     [0 − 1.0·0,   1 − 1.0·0 ] ]
  = [ [1, −0.44],
      [0,     1 ] ]

b = diag(p)·ex − F = [1.1·0.50 − 0.05,
                      1.0·1.10 − 0.10]
  = [0.55 − 0.05, 1.10 − 0.10]
  = [0.50, 1.00]
```

### 5.4. Solving for x

Since M is upper triangular:
1) Equation for x₂:
```
x₂ = b₂ / M₂₂ = 1.00 / 1 = 1.00
```
2) Equation for x₁:
```
x₁ − 0.44·x₂ = 0.50  ⇒  x₁ = 0.50 + 0.44·1.00 = 0.94
```

So the relative end-of-day wealth vector is:
```
x = [0.94, 1.00]
```

### 5.5. Fractional Yield and Interpretation

Fractional gain = x − 1:
```
Gain = [0.94 − 1.00, 1.00 − 1.00] = [−0.06, 0.00]
      = [−6%, 0%]
```
Bank A loses 6%, and B breaks even relative to its starting capital.

### 5.6. Scaling to Actual Capital

If we start with k = [10 000, 100], then actual end-of-day wealth:
```
wealth = k ⊙ x = [10 000·0.94, 100·1.00] = [9 400, 100]
```
Absolute gain/loss:
```
Δwealth = wealth − k = [9 400 − 10 000, 100 − 100] = [−600, 0]
```
