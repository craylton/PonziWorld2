
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

HOW TO RESOLVE THIS:
EXAMPLE 1
So imagine we have a situation where A is invested in B, and B is invested in the stock market
We denote this as follows:
    A   B   S
A | 0   a_b 0   |
B | 0   0   b_s |

This means A has invested 0 into itself (as that isn't allowed), has invested `a_b` into B, and 0 into the stock market (S)
Meanwhile B has invested `b_s` into the stock market, and not invested anywhere else

Now there are other factors at play here
    A and B have what I'm calling the 'ponzi factor `P-a` and `P_b` - how much they are lying about their profit
    A and B have to pay their own fixed costs `f_a` and `f_b`
    The stock market changes by a certain yield, y_s

So B started with a total of `b_s` invested, and has now gained money from stocks, but also paid fixed costs
    So B now has `b_s * y_s` invested in stocks
    B's overall profit is therefore `b_s * y_s / b_s = y_s` in this simple case
    And then we simply apply the fixed costs to figure out how much money B now has
    Worth noting that the fixed costs always come out of cash
        in this example both players start with 0 cash
        which means B's cash will be `-f_b` after paying their fixed costs
    so overall B started with a capital of `b_s`, and now has capital of `b_s * y_s - f_b`
    And finally we can multiply by the ponzi factor to figure out what B 'claims' to have
    So `B_final = (b_s * y_s - f_b) * p_b`
Now we can consider A. A is only invested in B, so we can calculate from there
    B started with `B_start = b_s`, so their yield was `y_b = B_final / B_start`
    That means A's investment into B of `a_b` becomes `a_b * y_b`
    And since that's the only investment A had, we can figure out `A_final`
    Following the same logic as with B: `A_final = (a_b * y_b - f_a) * p_a`

One small thing worth mentioning:
    I started by calculating B, and then calculated A afterwards
    This is because B had no investemnts into other banks - only into stocks where we already know the yield
    That's a little trick - you can only calculate a bank's yield once you know the yields of all its investments

That's a pretty simple example. Let's look at a more complex one now

EXAMPLE 2

    A   B   C   D   Sto Cash
A | 0   a_b 0   a_d 0   a_cash  |
B | 0   0   b_c 0   0   b_cash  |
C | c_a c_b 0   0   c_s c_cash  |
D | 0   d_b 0   0   0   d_cash  |

And again the banks all have their own ponzi factors and fixed costs.
And again, we know the yield of stocks (`y_s`) and the yield of cash (always 1, cash never changes).

In this exmple we have a few loops:
A -> B -> C -> A
A -> D -> B -> C -> A
B -> C -> B

The way we tackle loops is by reducing them down and treating them as cash
To explain this, let's look at the B -> C -> B loop
if, for example `b_c = 50`, and `c_b = 30`, we treat it as though `b_c = 20`, and `c_b = 0`
and for both of those banks, we treat it as though the remaining 30 is in cash
so in that example we would reduce the matrix down to:

    A   B   C   D   Sto Cash
A | 0   a_b 0   a_d 0   a_cash      |
B | 0   0   20 0    0   b_cash + 30 |
C | c_a 0   0   0   c_s c_cash + 30 |
D | 0   d_b 0   0   0   d_cash      |

If we repeat this reduction for each loop, eventually we'll end up with a bank which is only invested in known yields
In this case the only known yields are stocks and cash
Looking at the matrix, none of our banks fit that requirement
this means we need to find another loop that needs to be reduced

Let's reduce the A -> B -> C -> A loop next
if we say a_b = 20 and c_a = 40
after reducing, the new matrix will look like this:

    A   B   C   D   Sto Cash
A | 0   0   0   a_d 0   a_cash + 20 |
B | 0   0   0   0   0   b_cash + 50 |
C | 20  0   0   0   c_s c_cash + 50 |
D | 0   d_b 0   0   0   d_cash      |

B is now only invested in known yields (cash + stocks), so we can calculate B's yield is p_b
And now D is only invested in known yields (cash + B), so D's yield is `p_b * p_d`
And now A (invested in D and cash) has a yield of `p_b * p_d * p_a`
so finally, C started with `20 + c_s + c_cash + 50`
C ends up with `B_final = 20 * p_b * p_d * p_a + c_s * y_s + c_cash + 50`
And then multiply that by the ponzi factor p_c to get the yield

This solution does have a problem though
if we look back at the original matrix, B had invested 50 in C
at the end of the calculations, C has made a positive yield `y_b > 1`
but B will be feeling a bit sad to see their investemnt of 50 remains at 50
'what gives? I invested 50, it went up, but I still ended up with 50?'

So ultimately this solution solves the problem of people creating infinite money by reinvesting into each other
But breaks down when matching up percentage yields to actual profits