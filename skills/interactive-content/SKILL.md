---
name: interactive-content
description: Build interactive visualizations, charts, plots, and simulations inside sitegen knowledge bases using template shortcodes — including custom canvas sketches with arbitrary JavaScript. Use when a topic would benefit from something the reader can watch, manipulate, or test themselves.
---

# Interactive content in sitegen

The whole point of sitegen is that static pages can still be *alive*. Reach for
an interactive element whenever a concept is dynamic, geometric, or has a
"feel" that prose can't convey. Always run `sitegen templates` to confirm
argument names; missing required arguments are build ERRORs.

## Charts — data you have

```markdown
{{< chart type="line" data="0,1,4,9,16,25" title="x²" ylabel="energy" >}}
{{< chart type="bar"  data="3,7,2,9,4" color="#8b7cff" >}}
{{< chart type="scatter" data="0:1 1:3 2:2 3:5" >}}
```

`data` is comma-separated y-values, or `x:y` pairs (space/semicolon separated),
or JSON. Types: `line` (filled), `bar`, `scatter`.

## Plots — functions you can express

```markdown
{{< plot fn="Math.sin(x)" title="sine" >}}
{{< plot fn="x*x ;; Math.exp(x)-1" xmin="-2" xmax="2" ymin="-1" ymax="6" >}}
```

`fn` is a JavaScript expression in `x` (all `Math.*` functions are in scope).
Separate multiple curves with `;;`.

## Built-in simulations — physics out of the box

```markdown
{{< sim name="pendulum"        length="1.5" angle="2.4" damping="0.05" >}}
{{< sim name="doublependulum"  angle1="2.4" angle2="2.6" >}}
{{< sim name="lorenz"          sigma="10" rho="28" beta="2.667" >}}
{{< sim name="nbody"           bodies="6" g="1" >}}
{{< sim name="life"            cell="10" >}}
{{< sim name="vectorfield"     fx="y" fy="-x" >}}
```

Each renders to an animated canvas with play/pause/reset controls. `name` is
required; all other arguments are optional configuration (the `sim` template
accepts arbitrary extra args, so tune freely). These are ideal for
dynamical-systems topics: attractors, chaos, oscillations, gravitation.

## Custom sketches — anything you can draw

When no built-in fits, write the visualization yourself with `{{< sketch >}}`.
The inner block is JavaScript executed **every animation frame** with these
locals in scope:

| local    | meaning                                            |
|----------|----------------------------------------------------|
| `ctx`    | 2D canvas context (already cleared each frame)     |
| `W`, `H` | canvas size in CSS pixels                          |
| `t`      | seconds since start                                |
| `dt`     | seconds since last frame                           |
| `frame`  | frame counter (use `if (frame===0){…}` for setup)  |
| `mouse`  | `{x, y, down}` pointer state                       |
| `state`  | a persistent object to stash values between frames |

```markdown
{{< sketch height="320" caption="A traveling wave you can poke." >}}
ctx.strokeStyle = '#5b9cff';
ctx.lineWidth = 2;
ctx.beginPath();
for (let x = 0; x <= W; x++) {
  const phase = (x / W) * 12 - t * 2;
  const drive = mouse.down ? (1 - Math.abs(x - mouse.x) / 80) : 0;
  const y = H/2 + Math.sin(phase) * 40 * (1 + Math.max(0, drive));
  x === 0 ? ctx.moveTo(x, y) : ctx.lineTo(x, y);
}
ctx.stroke();
{{< /sketch >}}
```

Rules for good sketches:
- Draw relative to `W`/`H` so it scales; never hardcode pixel sizes.
- Read colors from CSS variables for theme-consistency when you can, e.g.
  `getComputedStyle(document.documentElement).getPropertyValue('--accent')`.
- Use `state` for simulation variables; initialize them on `frame===0`.
- Keep per-frame work bounded — this runs at 60fps.
- Invite interaction: respond to `mouse`. A reader who can *poke* it learns more.

## Self-check questions

```markdown
{{< quiz question="Doubling a pendulum's length changes its period by…"
         options="½× | √2× | 2× | 4×"
         answer="2"
         explain="T ∝ √L, so length ×2 gives period ×√2." >}}
```

## Layout helpers

```markdown
{{< columns count="2" >}}
Left column markdown… (including [[links]] and **formatting**)
{{< /columns >}}

{{< figure src="diagram.svg" caption="Phase portrait." >}}
{{< eq number="3" >}}\oint \vec{B}\cdot d\vec{l} = \mu_0 I{{< /eq >}}
```

## Workflow

1. Add the shortcode to your markdown.
2. `sitegen check <dir>` — fixes argument errors instantly.
3. `sitegen build <dir> -o _site` and open the page to watch it run.
4. Iterate on the visual until it genuinely illuminates the idea.
