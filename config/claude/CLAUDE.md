## General

Do not tell me I am right all the time. Be critical. We're equals. Try to be neutral and objective.

Do not excessively use emojis.

Prefer using browser agent skill over using playwright directly.

## Writing docs / README
Never use dashes (— or -) as punctuation in documentation or README files. Rephrase sentences using periods, commas, or parentheses instead.

## Coding Standards
When working with Laravel/PHP projects, always use the php-guidelines-from-spatie skill

Avoid `private const` in PHP. Replace each one: inline it if it is used once (use a descriptive local variable when a raw literal would obscure the meaning the name carried), or turn it into a private property if it is used more than once (`private static` when the using methods are static). When a constant is referenced from a default parameter value or a PHP attribute, where no property can be used, inline the literal there.

## Using GitHub
For questions about GitHub, use the gh tool
Never mention Claude Code in PR descriptions, PR comments, or issue comments
Do not include a "Test plan" section in PR descriptions
Keep PR descriptions terse. No section titles. Focus on the main things; minimal examples are fine. Do this unless I ask for more detail.
