---
title: Driven Oscillator
tags: [oscillations]
summary: Push a damped oscillator with a periodic force and, after the transients fade, it forgets its own rhythm and dances to yours.
weight: 24
---

# Driven Oscillator

Left alone, a [[Damped Oscillator]] eventually stops. Keep pushing it with a periodic external force and it settles instead into a steady rhythm set by the *drive*, not by itself. With a sinusoidal forcing at angular frequency $\omega$, the equation gains a right-hand side:

$$m\ddot x + c\dot x + kx = F_0\cos(\omega t).$$

This is one of the most consequential equations in physics and engineering — it governs everything from a child on a swing to the tuning circuit in a radio.

## Transient plus steady state

The solution splits cleanly into two pieces:

$$x(t) = \underbrace{x_h(t)}_{\text{transient}} + \underbrace{x_p(t)}_{\text{steady state}}.$$

- The **transient** $x_h$ is the damped oscillator's own response to being disturbed. It decays as $e^{-\zeta\omega_0 t}$ and is gone after a few time constants — the system *forgetting* its initial conditions.
- The **steady state** $x_p = A(\omega)\cos(\omega t - \delta)$ oscillates forever at the *driving* frequency, lagging it by a phase $\delta$.

{{< chart type="line" data="0,0.9,0.3,1.4,0.7,1.7,1.0,1.85,1.2,1.9" title="Buildup to steady state" caption="The response grows from rest as the transient fades, levelling off into a constant-amplitude driven oscillation." >}}

## Amplitude depends on tuning

The steady-state amplitude

$$A(\omega) = \frac{F_0/m}{\sqrt{(\omega_0^2-\omega^2)^2 + (c\omega/m)^2}}$$

is small when you drive far above or below the natural frequency $\omega_0$, and it swells dramatically when $\omega$ approaches $\omega_0$. That peak is [[Resonance]] — important enough to deserve its own page. Push the system harder, into the nonlinear regime, and steady oscillation can detach from the drive entirely, becoming a self-sustained [[Limit Cycle]].

{{< note kind="key" title="The system has no memory — eventually" >}}
The defining feature of a driven *linear* oscillator is that its long-term behavior is independent of how it started. Two identical oscillators released differently converge to the *same* steady oscillation once their transients die. The drive wins; the initial conditions are forgotten.
{{< /note >}}

## See also

- [[Resonance]]
- [[Limit Cycle]]
- [[Damped Oscillator]]
