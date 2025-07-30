We have a collection of banks all invested in one another, as well as having investments in exogenous assets too. We define the following variables:

n = the number of banks.
m = management fees. This is a constant, approximately 0.005. If bank A invests £1000 in bank B, bank B will earn £5 as cash from that investment at the end of the day (regardless of how well or badly the money was invested).
e = exogenous investments. This is an n*x array which represents how much each bank has invested into each asset. So for example, if bank A has invested 15 into stocks but has no cash, and bank B has invested 7 into stocks and has 5 cash, e will look like [[0, 15], [5, 7]].
y = exogenous asset yields. This n*1 vector represents how well each exogenous asset has performed, as a percentage. So if stocks went up by 10%, but cash didn't change (cash never changes), then y will look like [1, 1.1].
W = bank investments. Similar to the exogenous investments array, this is an n*n array that represents how much each bank has invested into each other bank. So if A has invested 3 into bank B, and bank B has invested 9 into bank A, W will look like [[0, 3], [9, 0]]. Note that banks cannot invest into themselves.
c = total amount capital managed by each bank. You can calculate c by adding up each row in W, plus each row in e. So for the examples given above, c will be [15 + 3, 12 + 9] = [18, 21].
p = ponzi factor. This is an n*1 vector that represents how much each bank is willing to lie about their performance. So while bank A might only make a 2% gain, they can lie to make it appear that they actually made 5%, for example.
F = fixed costs. Each bank has their own fixed costs that they must pay each day. this is another n*1 vector.

Define:
- 1ₙ ∈ ℝⁿ : column vector of ones  
- 1ₓ ∈ ℝˣ : column vector of ones  
- c = W·1ₙ + e·1ₓ ∈ ℝⁿ : total capital managed by each bank  
- s = Wᵀ·1ₙ ∈ ℝⁿ : total funds invested into each bank by others  
- m ∈ ℝ : management‐fee rate (scalar)  
- p, F ∈ ℝⁿ : ponzi factors and fixed costs  

Compute:  
Δ_exo = e·(y − 1)           ∈ ℝⁿ  
g      = (Δ ⊘ c) + p       ∈ ℝⁿ  
fee    = m·s               ∈ ℝⁿ          (fee income on in-flows)  
inter  = W·(g − m)         ∈ ℝⁿ          (interbank gains net of fees)  

Total net payoff (implicit in Δ):  
Δ = Δ_exo + inter + fee − F  

⇔  (Iₙ − W·diag(1 ⊘ c))·Δ = Δ_exo + W·(p − m) + m·s − F  

Updated investment matrices:  
e₁ = e · diag(y)               (scale each exogenous column by its yield)  
Let α = g − m + 1 ∈ ℝⁿ  
W₁ = W·diag(α)                (so W₁ᵢⱼ = Wᵢⱼ·αⱼ)  

That's basically the whole algorithm. There are a couple of small tweaks that need to be made however.
Taxes:
If bank A invests 100 in bank B, and then bank B invests that bank into bank A, and finally bank A puts it into stocks - you can see how bank A's success causes bank B's success, which in turn causes more success for bank A (since A is invested in B). As a result of this feedback loop, it's possible that we end up with more money in the system than we should.
To counter this, we use 'taxes'. We can detect how much money has been 'invented' in the system, and simply add that amount to the banks' fixed costs to recoup the fictional money.
Note that the heuristcs for calculating how much taxes to charge each individual bank has not yet been determined. Clearly some banks contribute more to the problem than others (e.g. those heavily invested in loops vs not in loops), but it is a difficult problem to figure out how much.
Debt:
In order to avoid 'divide-by-zero' errors, and bizarre negative values where we don't expect them, we need to ensure that every investment is a positive value. However, we also need to include cash in these calculation, while cash can easily be zero or negative.
For example, if a bank invests 100 into an asset and turns it into 150, you might say the bank knows what it's doing. But if you learn that bank is also holding 1,000,000 in cash, completely uninvested, then you'd say it's a pertty terrible bank only managing to gain 50 from over 1m in capital. Which is why it's important to consider cash.
For these reasons, for the sake of these calculations at least, we restrict cash from being negative. If cash ever does go negative, we simply set it to 0, and perform these calculations as though the cash were zero. The debt will be processed in another step (hefty interest payments).

-----------------------------------
SIMPLE EXAMPLE

A invests 100 in B, 300 in stocks, and holds no cash. B invests 200 in A, holds 100 in cash. Stocks then go up 10%.

    e = [[0, 300], [100, 0]]
    y = [1, 1.1]
    W = [[0, 100], [200, 0]]
    m = 0.005
    p = [0.01, 0]
    F = [10, 8]
    c = [300+100, 100+200] = [400, 300]

Finding Δ:

    Δ_a =
    0 * 0 + 300 * 0.1 +
    0(Δ_a/400 + 0.01 - 0.005) + 100(Δ_b/300 + 0 - 0.005) +
    (0 + 200) * 0.005 -
    10

    Δ_a = 0 + 30 + 0 + (Δ_b/3 - 1/2) + 1 - 10
    = 41/2 + Δ_b/3

    Δ_b = 
    100 * 0 + 0 * 0.1 +
    200(Δ_a/400 + 0.01 - 0.005) + 0(Δ_b/300 + 0 - 0.005) +
    (100 + 0) * 0.005 -
    8

    Δ_b = 0 + 0 + (Δ_a/2 + 1) + 0 + 1/2 - 8
        = Δ_a/2 - 13/2

    Δ_a = 41/2 + Δ_b/3
    Δ_b = Δ_a/2 - 13/2

Solving these simultaneous equations gives:

    Δ_a = 22
    Δ_b = 4.5

From this we can calculate the new investment sizes:

    e_1 = [[0, 330], [100, 0]]
    W_1 = W((Δ/c + p - m) + 1)
    W_1 = [[0((22/400 + 0.01 - 0.005) + 1), 100((4.5/300 + 0 - 0.005) + 1)],
        [200((22/400 + 0.01 - 0.005) + 1), 0((4.5/300 + 0 - 0.005) + 1)]]
        = [[0, 101], [212, 0]]

**Expected vs. Actual System Net Worth**

1. Gross system assets:
   - Starting total: 400
   - Stock gain: +30
   - Ponzi gain (A): +2  
   → **Gross total = 432**

2. Fixed costs paid:
   - A pays: 10
   - B pays: 8  
   → **Total fixed costs = 18**

3. Expected net total:
   432 − 18 = **414**

4. Compare exogenous plus interbank balances: 422 − 212 + 304.5 − 101 = 413.5.

**Result:** there is a 0.5 shortfall relative to the expected 414, so the system would owe a 0.5 ‘tax rebate’ to the banks.

-----------------------------------------------

## Example Using Bank Yield Instead of Δ

We define each bank's yield:

    y_A = c_1_A / 400
    y_B = c_1_B / 300

where c_1 = e·y_exo + W·(y_bank + p - m) - F, and y_bank = [y_A, y_B].

Compute each term:

- e·y_exo for A: 0×1 + 300×1.1 = 330
- F_A = 10

- e·y_exo for B: 100×1 + 0×1.1 = 100
- F_B = 8

Then the W·(y_bank - m) terms:

- For A: 100·(y_B - 0.005)
- For B: 200·(y_A + 0.01 - 0.005)

So we get:

    c_1_A = 330 + 100*(y_B - 0.005) - 10
          = 320 + 100*(y_B - 0.005)

    c_1_B = 100 + 200*(y_A + 0.01 - 0.005) - 8
          = 92 + 200*(y_A + 0.005)

Yield equations:

    y_A = (320 + 100*(y_B - 0.005)) / 400
    y_B = (92 + 200*(y_A + 0.005)) / 300

Solving:

    y_A = 1.0515
    y_B = 1.011

Then after subtracting the fixed costs from cash, the updated investment matrices are:

    e₁ = [[-10, 330], [92, 0]]
    W₁ ≈ [[0, 100.6], [211.3, 0]]

So A's net worth has gone from 200 to 209.3. B's net worth went from 200 to 202.7. So overall 12 has been added to the system. Which is exactly what we expect, since 10 + 8 was paid in fixed costs and 30 was gained from stocks: 30-18 = 12.