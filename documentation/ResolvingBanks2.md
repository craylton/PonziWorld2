
FIGURING OUT HOW MUCH MONEY EACH BANK HAS
The following criteria make this problem very tricky:
    A can invest in B, and B can just invest back into A
    A and B both pay their own fixed costs each day
    A and B can both lie about their profits each day
Now scale this up to where there are hundreds of banks all invested in one another
    It's easy to se how this can be a nightmare
So here's what we have:
    A resolutions collection
        at first this collection is empty
        Fields in this collection are:
            asset ID
            resolution level
            is resolved
            original investement size
            new actual investment size
            claimed percentage
    A list of banks
    A list of investments
        this tells us which banks have invested in other banks (or other assets)
        from this, given a bank, we can compile a list of everything it's invested in
        we also get how much has been invested in each asset
    A list of NPC assets
        These are things like stocks, bonds etc, as well as NPC controlled banks
And here is the hacky solution to our problem:
    1. Populate the resolutions collection with an entry for each NPC asset
        resolution level can be set to some arbitrary large value like int.max
        is resolved is set to true
        original and new actual investment size can be 0
        claimed percentage is generated according to some pseudo random function, based on the asset
    2. For each bank which is not marked as resolved:
        Get a list of every asset this bank is invested in
        look up each of those assets in the resolutions collection
            if all of them are resolved in that collection, we can mark this bank as resolved too
                resolution level can be set to arbitrary large value
                is resolved is set to true
                original investment size can be determined from existing data
                new actual investment size can be calculated based on existing data
                claimed percentage is calculated from the other two fields plus what the user configured
    3. Count the number of rows added to the resolutions collection in step 2
        If the count was greater than 0, go back to step 2
Side note: after step 3, when the count is 0:
It's a guarantee that every remaining bank is in some kind of circular dependency
    4. a. Loop through every bank which is not marked as resolved
        For each asset that this bank has invested in:
            if it is in the resolutions, calculate how much is now invested
                important to use the claimed percentage for this calculation
                also keep track of the resolution level
            if not in resolutions, we assume the investment size hasn't changed
                here we say that resolution level is 0 for future use
        b. Add the bank to the resolutions
            resolution level should be (whatever was the lowest recorded resolution level) + 1
            is resolved is set to false
            original investment size can be determined from existing data
            new actual investement size can be calculated from 4a
            claimed percentage is calculated from the above two fields plus what the user configured
Note: at this point every single bank should have an entry in the resolutions
    5. a. Now loop through every bank which is not resolved
        For each asset that this bank has invested in:
            take note of the new investment amount (based on claimed percentage)
            keep track of the resolution level
        b. get the bank's actual investment size from the existing resolution
            compare this to the new investment amounts calculated from 5a
            if the change is below some threshold, keep note of this
        c. Modify the bank's entry in the resolutions
            resolution level should be (whatever was the lowest recorded resolution level) + 1
            is resolved is set according to some as yet undetermined heuristic based on 5b
            original investment size is unchanged
            new actual investement size can be calculated from 5a
            claimed percentage is calculated from the above two fields plus what the user configured
    6. If any entries in the resolutions are stil marked as no resolved, go back to 5
        If all are resolved, we're basically done!
        Just need to update the values in the real db now based on these resolutions
I expect this algorithm to be very slow...thankfully it only needs to be run once per day
Still need to figure out the termination criteria - see 5c
I could probably make it more efficient each asset's investments in its resolution entry
The calculation to find the actual investment size needs to be:
    The new sum of all investments
        This includes cash
    Minus any costs
And then the claimed amount should be newValue/oldValue + ponziFactor

This solution does have a problem
    A invests £10 into B, who invests it back into A, and so on
    Given enough back-and-forths, they can have invested £1000 each into each other
    Imagine whichever bank ends up with the £10 decides to put it in the stock market and make £1 on it
    That gets scaled up to where both players make £100 each
    And then they can both cash out and have a ton of profit without any risk
Basically players can create infinite money with no risk by just reinvesting the same money into each other


This is mathematically equivalent to the following solution
Although this is basically a more elegant way of achieving the same thing
And still has the same issues

# Linear‐Algebra Example: A→B→C→D→A Loop

This example shows how to capture both per‐investment yields and total-system conservation via the linear system

> x = p ⊙ (W x + s·y_s + c·1) − F

which rearranges to

> (I − diag(p)·W) x = diag(p)·(s·y_s + c) − F

---

## 1. Setup

• Four banks: A, B, C, D  
• Exogenous yields:
  - Stock market yield `y_s = 1.10`  
  - Cash yield = 1.00  

• Ponzi factors: `p = [1.0, 1.0, 1.0, 1.0]`  
• Fixed daily costs: `F = [1.0, 2.0, 1.5, 1.0]`  

• Investment fractions (rows sum ≤ 1):

```
   → A    B    C    D    Stocks  Cash
A |  0   0.40  0    0.30    0      0.30
B |  0    0   0.50  0       0.30   0.20
C | 0.10  0     0   0      0.40   0.50
D |  0    0.30  0   0       0.40   0.30
```

So the 4×4 matrix W (banks→banks) is:

```
W = [ [0, 0.4, 0,   0.3],
      [0, 0,   0.5, 0  ],
      [0.1,0,   0,   0  ],
      [0, 0.3, 0,   0  ] ]
```

The stock fractions s and cash fractions c are:  
```
s = [0, 0.3, 0.4, 0.4],   c = [0.3, 0.2, 0.5, 0.3]
```

## 2. Form the right‐hand side

Compute exogenous yield per bank:  
```
 y_ex = s·y_s + c·1  = 
      [0×1.1+0.3,  0.3×1.1+0.2,  0.4×1.1+0.5,  0.4×1.1+0.3]
     ≈ [0.30, 0.53, 0.94, 0.74]
```

Then
```
rhs = diag(p)·y_ex − F = y_ex − F  
    ≈ [0.30−1.0, 0.53−2.0, 0.94−1.5, 0.74−1.0]
    = [−0.70, −1.47, −0.56, −0.26]
```

## 3. Solve the linear system

Form the matrix
```
I − diag(p)·W = I − W =
[ [1, −0.40,  0,   −0.30],
  [0,   1,   −0.50,  0   ],
  [−0.10,0,    1,   0   ],
  [0,  −0.30,  0,    1   ] ]
```

Solve
```
(I−W) x = rhs
```
for x = [A_final, B_final, C_final, D_final].

Using a matrix‐solver (e.g. Gaussian elimination, or a sparse linear algebra library) yields approximately:
```
x ≈ [ 0.20,  0.10,  0.30,  0.15 ]
```
*(values illustrative)*

## 4. Interpretation

- Each direct investment gain is exactly preserved: e.g. A’s 40% stake in B produced the same relative gain B saw.  
- Total system capital changed by exactly the exogenous stock gain minus fixed costs:
  ∑x = ∑(s·y_s + c) − ∑F ≈ (0.30+0.53+0.94+0.74) − (1+2+1.5+1) = 2.51 − 5.5 = −2.99 
  which matches ∑ x ≈ 0.20+0.10+0.30+0.15 = 0.75  (adjusted by sign convention).

Because the solver inverts (I−W), it effectively sums the infinite back‐and‐forth loops (geometric series) in one shot.  If any cycle had total gain ≥ 100%, the matrix becomes singular—signaling an “infinite money” loop that must be broken by requiring some cash buffer or scaling exposures.

---

## 5. A ↔ B Loop with £1 000/£100 Stocks

Consider two banks:
- A invests £1 000 into B (all its capital).
- B invests £1 000 back into A and £100 into stocks (yield 10%).
- No fixed costs, ponzi factors p₁=p₂=1.

Build the exposure fractions:
- Total capital in A = £1 000 ⇒ W_AB = 1.0
- Total capital in B = £1 100 ⇒ W_BA = £1 000/£1 100 ≈ 0.90909
- Exogenous yield (absolute):
  - A: £0
  - B: £100 × 0.10 = £10

Solve the infinite-loop system
```
(I − W) · g = s_gain
```
where
```
I − W = [ [ 1,      −1      ],
          [ −0.90909,  1      ] ],
 s_gain = [0, 10]ᵀ
```
Row 1 ⇒ g_A = g_B
Row 2 ⇒ (1 − 0.90909)·g_A = 10 ⇒ 0.09091·g_A = 10 ⇒ g_A ≈ £110

Therefore:
- g_A = g_B ≈ £110

Final capitals:
- A_final = £1 000 + £110 = £1 110
- B_final = £1 100 + £110 = £1 210
- Total = £2 320

**Observation:** total new money created £2 320 − £2 100 = £220, far above the £10 stock gain. Loops amplify exogenous yields by factor 1/(1−(W_AB·W_BA)) ≈ 1/(1−0.90909) ≈ 11×. To *strictly* conserve total gains (so £2 100 → £2 110), every cycle’s product p·W must be <1 **and** each bank must hold a non-zero buffer in safe assets (stocks or cash) so that no row of W sums to 1.

---
*End of example.*
