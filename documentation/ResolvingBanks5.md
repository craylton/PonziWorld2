We have a collection of banks all invested in one another, as well as having investments in exogenous assets too. We define the following variables:

n = the number of banks.
m = management fees. This is a constant, approximately 0.005. If bank A invests £1000 in bank B, bank B will earn £5 as cash from that investment at the end of the day (regardless of how well or badly the money was invested).
e = exogenous investments. This is an n*x array which represents how much each bank has invested into each asset. So for example, if bank A has invested 15 into stocks but has no cash, and bank B has invested 7 into stocks and has 5 cash, e will look like [[0, 15], [5, 7]].
y = exogenous asset yields. This n*1 vector represents how well each exogenous asset has performed, as a percentage. So if stocks went up by 10%, but cash didn't change (cash never changes), then y will look like [1, 1.1].
W = bank investments. Similar to the exogenous investments array, this is an n*n array that represents how much each bank has invested into each other bank. So if A has invested 3 into bank B, and bank B has invested 9 into bank A, W will look like [[0, 3], [9, 0]]. Note that banks cannot invest into themselves.
c = total amount capital managed by each bank. You can calculate c by adding up each row in W, plus each row in e. So for the examples given above, c will be [15 + 3, 12 + 9] = [18, 21].
p = ponzi factor. This is an n*1 vector that represents how much each bank is willing to lie about their performance. So while bank A might only make a 2% gain, they can lie to make it appear that they actually made 5%, for example.
F = fixed costs. Each bank has their own fixed costs that they must pay each day. this is another n*1 vector.

The core update equation is (all vectors are n×1, matrices are sized appropriately):
Δ = e · (y − 1)                    // exogenous asset gains (n×1)
    + W · ((Δ ⊘ c) + p − m)        // interbank investment gains (n×1)
    + W m                          // cash from management fees (n×1)
    − F                            // fixed costs (n×1)

Here:
- Δ, p, m, F, c are n×1 column vectors.
- e is n×x, y is x×1, so e·(y−1) yields n×1.
- W is n×n, multiplication produces n×1.
- ⊘ denotes element-wise division of two n×1 vectors (Δ ⊘ c = [Δ_i / c_i]).
- W m multiplies W by scalar m, giving an n×1 vector of fee income.

Then the updated investment matrices are:
e₁ = e ∘ diag(y)                // element-wise scaling of exogenous assets
W₁ = W · diag((Δ ⊘ c) + p − m + 1)  // adjust interbank weights by realized returns

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
    = Δ_a/2 - 9/2

Δ_a = 41/2 + Δ_b/3
Δ_b = Δ_a/2 - 9/2

Solving these simultaneous equations gives:

Δ_a = 22.8
Δ_b = 6.9

From this we can calculate the new investment sizes:

e_1 = [[0, 330], [100, 0]]
W_1 = W((Δ/c + p - m) + 1)
W_1 = [[0((22.8/400 + 0.01 - 0.005) + 1), 100((6.9/300 + 0 - 0.005) + 1)],
    [200((22.8/400 + 0.01 - 0.005) + 1)], [0((6.9/300 + 0 - 0.005) + 1)]]
    = [[0, 101.8], [212.4, 0]]

So we end up with A having 101.8 invested in B, 330 invested in stocks, and 200 * 0.005 = 1 in cash (from the management fee).
And B ends up with 212.4 invested in A, and 100 + 100*0.005 = 100.5 in cash.

You can see from this that the total amount of money in the system used to be 400. Now 18 has been paid as fixed costs, 30 was generated from stocks, and also A 'invented' 2 (200 * 0.01) via their ponzi factor. So you'd expect the net worth of the system to now be 450.

However after all this we end up with:
432.8 - 212.4 + 312.9 - 101.8 = 431.5
So actually in this case we're missing 18.5, which means we'd need to give a 'tax rebate' to the banks.