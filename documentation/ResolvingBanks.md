

FIGURING OUT HOW MUCH MONEY EACH BANK HAS
The following criteria make this problem very tricky:
    A can invest in B, and B can just invest back into A
    A and B both pay their own fixed costs each day
    A and B can both lie about their profits each day
Consider this scenario:
    A invests £10 into B, who invests it back into A, and so on
    Given enough back-and-forths, they can have invested £1000 each into each other
    Imagine whichever bank ends up with the £10 decides to put it in the stock market and make £1 on it
    That gets scaled up to where both players make £100 each
    And then they can both cash out and have a ton of profit without any risk
Consider a similar scenario:
    you can have a loop of multiple players who invest in each other
    Like A -> B -> C -> A
    Then the above problem applies but is trickier to spot
Now scale this up to where there are hundreds of banks all invested in one another
    It's easy to se how this can be a nightmare


I have two potential solutions to this, outlined in the files in this folder
However both of them have issues when it comes to loops

Let's look at an example scenario:
    Bank A invests £100 into bank B, bank B reinvests all £100 back into A, and then A invests that £100 into stocks. We now have A and B both having £100 in one another, plus A has £100 in stocks.
    Now stocks go up 10%. How do we resolve this system?
    Well, technically there is only £100 in this system. A thinks it has £200, and B thinks it has £100, but actually it's all just the same £100 being reinvested.
    
    So the first requirement is that the total amount of money in this system must go up by exactly £10 (10% of 100).
    If we ignore the first requirement, then banks are able to conjure up free money out of thin air, simply by reinvesting the same money back and forth. This is why in this example the total money in the system must go up by exactly £10.
    
    Now consider we implement this in the simplest way we can - we give all of that £10 to bank A (since A is invested in stocks after all). After doing this, we see A thinks it started with £200 (100 in B and 100 in stocks). At the end it ends up with £210 (a 5% gain) after we give it the £10.
    Now B will see that A has gained 5%. Since B invested £100 into A, B will be expecting a 5% return on its £100 investment, or £5. But instead it gets nothing!

    So the second requirement is that the perceived yields must all match up. If a bank invests into some asset, and that asset increases by some amount, the bank should get a profit accordingly.
    If we ignore this second requirement, then the system will be broken - banks will assume there is a bug.

After looking through this example, I'm pretty sure there is no solution which satisfies both requirements.
So perhaps there is some middle ground? A compromise?
Perhaps I can just lie to the banks and to make it look like all the numbers add up when in reality they don't...although then I'd need to store a value for each pair of banks. So if there are 10000 banks, I'd have to store 100 million fake values per day - not an ideal solution.
What might be useful is that the amount of 'free money' that banks can conjure up is based on the proportion of money stored in 'known assets' like cash or stocks. So if I enforce that banks can't invest more than that, say 80% of their money into other banks, this would limit the impact of this.
Another option is to keep track of the amount of free money earnt by each bank (is this possible?), and to simply increase those banks' fixed costs accordingly. This is handy because the players who figure out and intentionally exploit this quirk won't gain any more than people who accidentally join an 'investment loop'. It does mean though that some players could end up with ludicrously high fixed costs without any explanation - I'd have to be careful how I present this to the player.

Another problem - probably more serious - doesn't even involve loops.
Imagine A invests £100 into B, and B invests it all into stocks.
Clearly there is only £100 in this system. So when stocks go up, the total amount of money in the system should go up by £10.
But who should get the £10?
B is the one who actually made the gain by putting it in stocks
But A is the one who actually owns the £100
Do we share the money £5 each? What if it's a chain of 10 people, do they all get £1 each? Why would they bother investing in each other if they're only getting a fraction of the actual gains - surely they should only bother investing in stocks and never bother with anything else.

We can make B charge a management fee to A
So that way A would get say £9 and B gets £1
Technically it's A's money so they should get most of the profit

-----------------------------------------

Alright, I've had a bigger think about this and have come up with the following.
Each bank has a list of exogenous assets that they have invested in, call it a.
Each bank has a list of other banks that it has invested in, call this list b.
And each bank has a list of other banks that have invested in it, call that c.
Now, the net worth of each bank before we do our calculations will be a + b - c (as in, add up all the things we are invested in and subtract all the things invested in us)
And after our calculations, the new net worth of that bank will be the sum of each asset in a multiplied by its yield, plus the sum of each bank in b multiplied by its yield multiplied by (1 - that bank's commission/management fee), minus the sum of each asset in c multiplied by this bank's yield multiplied by this bank's commission.
Or as an equation, if we're trying to calculate bank A's values:
(∑e_ai·y_ei)·m_a + ∑W_ai·(1-m_i) - ∑W_ia·y_a·m_a - F_a
Where W is a matrix of how much each bank has invested in each other.
So let's look at a simple example

A is invested in B, B is invested in stocks and also holds cash.
So W = [[0, 100], [0, 0]]
and e (previously called a) is [[0, 0], [10, 100]] (B has 10 cash and 100 stocks)
We'll say that bank A has fixed costs of 2, and B has fixed costs of 1. So F = [2, 1]
Finally, we'll say that A's ponzi factor (how much banks exaggerate their yield) is 1, and B's is 1.1. So p = [1, 1.1]. This means whatever yield B actually achieves, other banks will believe that they actually managed 10% more than that.
I'll also say that m, the management fee, is 0 if the bank's overall yield was negative, and 0.1 otherwise.

Now if we work through this example, we see that A started with a net worth of A_0 = 100, and ended with A_1 = 100·y_b·(1-m_b) - 2.
Meanwhile B started with B_0 = 100+10-100 = 10, and ended up with B_1 = 120 - 100·y_b·m_b - 1.
The overall (ponzi) yield y of each bank is y = [p_a·A_1/A_0, p_b·B_1/B_0]

So I essentially have some simultaneous equations to be solved.


-----------------------------

Again, working through that example revealed more problems. The yield of a bank is actually very difficult to work with. How do we even calculate the yield? Is it the rate of change of net worth? Or the rate of change of capital? In both cases, the 'before' value could be near zero - in which case the yield could be huge, or negative.
If your net worth was £-1 and you just made £10, that is clearly a good thing but your yield would be -1000% (I think). Which means despite making money, all your investors would lose 10 times more than they invested??

So clearly basing all the calculations on yield isn't going to work. I need a system that keeps track of how much money (absolute values) has been generated from various places, rather than multiplying everything by some factor.

So Introducing what I hope will be the final algorithm, outlined in ResolvingBanks5.md.