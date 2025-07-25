

# Summary of Bank Wealth Resolution Problem

This document provides an overview of the challenge of calculating daily bank wealth in a system with complex inter-bank investments, and it summarizes the various solutions that have been explored. The current proposed solution, which this document concludes with, is detailed in `ResolvingBanks5.md`.

## The Core Problem

Calculating the end-of-day wealth for a set of banks is complicated by several factors:
1.  **Investment Loops**: Banks can invest in each other, creating circular dependencies (e.g., A → B → C → A). These loops can artificially amplify gains, creating "free money" that violates the conservation of capital in the system.
2.  **Profit Distribution**: When one bank profits from an investment made with another's capital, it's unclear how to distribute those profits fairly while ensuring all parties perceive a correct yield on their investments.
3.  **Other Factors**: Daily fixed costs, management fees, and "Ponzi factors" (exaggerated profit reporting) add further complexity to the calculations.

A successful algorithm must satisfy two conflicting requirements:
- **Conservation of Capital**: The total money in the system should only change by the sum of external gains (e.g., from stocks) minus total costs. No money should be created or destroyed by internal transactions.
- **Consistent Yields**: If bank A invests in bank B, and bank B reports a 10% yield, bank A should see a 10% return on its investment in B.

## Explored Solutions

Several algorithms were devised to tackle this problem. While they provided valuable insights, each had drawbacks that led to further refinement.

1.  **Sequential Resolution (`ResolvingBanks1.md`)**: An intuitive approach where banks are resolved in dependency order. It identifies and "breaks" investment loops by converting the smallest part of the loop into a non-yielding asset (cash), allowing the rest of the system to be calculated.
    -   *Flaw*: This method is not fair, as the bank whose investment is used to break the loop does not receive the profit it should from that investment.

2.  **Iterative Resolution (`ResolvingBanks2.md`)**: A "hacky" but more robust iterative algorithm that repeatedly calculates bank wealth until the values converge. It handles loops by allowing their effects to propagate through the system over multiple iterations until a stable state is reached.
    -   *Flaw*: While it handles loops, it still allows for the creation of "infinite money" if not properly constrained.

3.  **Linear Algebra Approaches (`ResolvingBanks3.md`, `ResolvingBanks4.md`)**: These approaches model the problem as a system of linear equations, which can be solved simultaneously. This is a more elegant and mathematically sound way to handle the mutual dependencies. Different formulations were explored, one based on relative end-of-day wealth and another on absolute net worth.
    -   *Flaw*: These models also need mechanisms (like "taxes") to counteract the money-creation effect of high-yield loops and require careful formulation to be accurate.

## Current Proposed Solution

The insights from these earlier attempts led to the development of the algorithm detailed in **`ResolvingBanks5.md`**. This model uses a system of simultaneous equations to solve for the absolute change in each bank's holdings, providing what appears to be the most promising path to a correct and fair resolution of bank wealth.





# Sequential Resolution (ResolvingBanks1)

An initial approach to compute bank wealth by resolving dependencies and breaking loops.

Algorithm:
1. Mark exogenous assets (stocks, cash) as resolved.
2. Iteratively resolve any bank whose investments are entirely in resolved assets:
   - Calculate its end-of-day wealth using exogenous yields, fixed costs, and Ponzi factor.
   - Mark the bank as resolved.
3. If no banks can be resolved (all are in loops), assign a resolution level to each:
   - Compute tentative investments and claimed percentages based on existing resolutions.
   - Increase the level and repeat until values converge or meet a threshold.

Limitations:
- Converting looped investments to cash can misallocate profits.
- Performance degrades as the number of banks and loop depth increases.



# Iterative Resolution (ResolvingBanks2)

A multi-stage, collection-based approach to resolving bank wealth, followed by its linear-algebra reformulation.

## Hacky Multi-Pass Algorithm
1. Create a resolutions collection for exogenous (NPC) assets, marking them as resolved.
2. Iteratively process banks:
   - Resolve any bank whose investments depend only on resolved assets: compute its wealth and mark resolved.
   - When no further banks can be resolved, identify loops:
     - For each unresolved bank, compute tentative new investments and claimed yields.
     - Assign resolution levels to track loop depth.
   - Refine unresolved entries in subsequent passes until values converge or meet a threshold.

**Drawbacks:** Can be slow; loop heuristics and termination criteria require further tuning.

## Linear Algebra Equivalence
The iterative process can be expressed as a linear system:

```text
x = p ⊙ (W x + s·y_s + c) − F
(I − diag(p)·W) x = diag(p)·(s·y_s + c) − F
```
Solving this system handles interbank loops and preserves capital conservation in one step.



# Resolving Bank Wealth with Fees and Taxes (ResolvingBanks3)

A concise linear-algebra approach that incorporates management fees and correction of "free" loop gains via taxes.

## Notation and Inputs
- n banks; W (n×n) exposure fractions, s (n) exogenous asset fractions, c (n) cash fractions.
- y_s: exogenous asset yields, F (n): fixed costs, p (n): Ponzi factors, m: management fee rate.

## Algorithm
1. Compute exogenous returns: ex = s·y_s + c·1.
2. Adjust exposure for fees: A = W·(1−m).
3. Form and solve linear system for relative wealth x:
   (I − diag(p)·A)·x = diag(p)·ex − F.
4. If system is singular (infinite-loop gains), compute total loop surplus:
   G_loop = sum(x) − (sum(ex) − sum(F)),
   distribute taxes t (∑tᵢ = G_loop), and set x_final = x − t.

## Benefits and Notes
- Captures management fee retention in A.
- Solves interbank loops in one step.
- Taxes restore conservation when loops amplify gains.




# Net-Worth Resolution via Linear Algebra (ResolvingBanks4)

A step-by-step linear-algebra algorithm to compute end-of-day net worth, handling inter-bank loops, fees, and conservation.

## Notation
- n₀ᵢ: initial net worth of bank i; n₁ᵢ: end-of-day net worth.
- Wᵢⱼ: amount (or fraction) bank i invests in bank j.
- pᵢ: Ponzi factor (≥1), mᵢ: management fee rate, Fᵢ: fixed daily cost.
- Eᵢ: total exogenous returns (assets·yields).

## Core Equations
For each bank i:
```
n₁ᵢ = Eᵢ + ∑ⱼ[Wᵢⱼ·yⱼ·(1−mⱼ)] − ∑ₖ[Wₖᵢ·yᵢ·(1−mᵢ)] − Fᵢ
yᵢ  = pᵢ·(n₁ᵢ / n₀ᵢ)
```
This forms an N×N linear system in n₁.

## Solution Steps
1. Assemble LHS matrix M and RHS vector b from the above equations.
2. Solve M·n₁ = b using a linear solver.
3. If loops create “infinite” gains (singular M), impose a loop tax:
   - Compute surplus G_loop = ∑n₁ − (∑n₀ + exogenous_gain − ∑F).
   - Distribute tax t (∑tᵢ = G_loop) to restore capital conservation.

## Notes
- Eliminates iterative dependency resolution by inverting (I−diag(p)W).
- Preserves per-investment yields net of fees.
- Loop tax ensures no artificial money creation.
