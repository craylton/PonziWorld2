

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
What might be useful is that the amount of 'free money' that banks can conjure up is based on the proportion of money stored in 'known assets' like cash or stocks. So if I enforce that banks can't invest more that, say 80% of their money into other banks, this would limit the impact of this.
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

