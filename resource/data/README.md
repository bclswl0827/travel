# AK135 data provenance

The numerical travel-time data in this directory were transcribed from:

B. L. N. Kennett (2005), _Seismological Tables: ak135_,
Research School of Earth Sciences, Australian National University.

Official source:
https://download.iaspei.org/download/AK135tables.pdf

Copyright © B. L. N. Kennett 2005.

These files are machine-readable transcriptions of selected numerical tables
from the publication above, with the corrections and supplements below.
Numerical supplements and comparison values were calculated offline with
[TauP 3.2.1](https://github.com/crotwell/TauP/releases/tag/v3.2.1), using its
`ak135` model and a surface receiver.

## Deep S-wave supplement

The published PDF skips the deep S-wave table page between 25° and 50°.
The same gap exists in the ANU original archived on 2006-03-01:
https://web.archive.org/web/20060301122044/http://rses.anu.edu.au/seismology/AK135tables.pdf

The 26°–49° rows in `ak135_S_deep.txt` supplement that gap with the earliest
`S` arrival at source depths of 300–700 km in 50 km steps.
Travel times are rounded to 0.01 seconds and ray parameters to 0.01 seconds
per degree, matching the published table's precision. At 25° and 50°, the
calculated travel times differ from the published values by less than
0.02 seconds at all nine depths.

To reproduce the numerical supplement:

```sh
taup time --mod ak135 -h 300,350,400,450,500,550,600,650,700 \
  --deg 26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49 \
  -p S --first --json -o ak135_S_deep_supplement
```

The JSON `time` and `rayparam` fields provide the duration and slowness for
each `(distdeg, sourcedepth)` cell, respectively.

## Numerical corrections

The existing table grids were compared with reference calculations.
P and S use the earliest arrivals from `p/P/Pdiff` and `s/S/Sdiff`.
PKPab and PKPbc select the higher- and lower-ray-parameter
PKP branches, respectively. Other files use their named phase, with the
publication's PKPdf/SKSac/SKSdf labels corresponding to PKIKP/SKS/SKIKS.

Corrections and additions:

- P/S grids below 2° add 0.1° distance sampling at every existing source
  depth. This reduces the large interpolation errors of the original 1°
  grid near the source. Stray trailing page digits were removed from the
  P shallow 89° and P deep 69° rows.
- PcP at 97°/500 km is restored. Unavailable PcP cells at 97°/700 km and
  98°/500–700 km, and ScS at 98°/700 km, are all-zero cells rather than
  zero durations with nonzero ray parameters.
- PcP at 99°, ScP at 63°, and ScS at 99°–100° add all available arrivals
  at the existing source depths. Cells without arrivals remain all-zero.
- Invalid PKPbc at 155°/500–700 km, SKS at 144°/700 km, and SKIKS at
  105°/200 km are all-zero cells. SKS at 144°/500 km and SKIKS at
  105°/300–700 km use corrected arrival values. No runtime exclusion list
  is needed.
- The published SKP file combined SKP and SKIKP arrivals. `ak135_SKP.txt`
  now contains only the earliest SKP arrival, while `ak135_SKIKP.txt`
  contains the earliest SKIKP arrival. Both grids use 1° spacing and source
  depths 0, 35, 50, 100, 200, 300, 500, and 700 km. Their distance ranges
  are 129°–149° and 110°–180°.

To reproduce the reference grid for the audit:

```sh
taup time --mod ak135 \
  -h 0,15,35,50,100,150,200,250,300,350,400,450,500,550,600,650,700 \
  --degreerange 1 \
  -p p,P,Pdiff,s,S,Sdiff,PcP,ScP,ScS,PKP,PKIKP,SKP,SKIKP,SKS,SKIKS \
  --json -o ak135_reference

taup time --mod ak135 \
  -h 0,15,35,50,100,150,200,250,300,350,400,450,500,550,600,650,700 \
  --deg 0.1,0.2,0.3,0.4,0.5,0.6,0.7,0.8,0.9,1.0,1.1,1.2,1.3,1.4,1.5,1.6,1.7,1.8,1.9 \
  -p p,P,s,S --json -o ak135_near_source
```

At the stored grid points, the retained publication values differ from
this reference by less than 0.06 seconds; generated values are rounded to
0.01 seconds. This is a comparison of grid data, not an error bound for
interpolated distances and depths. Interpolation still requires valid
adjacent cells and does not extrapolate through phase-existence gaps.

## Licensing

The Apache License 2.0 for this repository applies to the original software
source code. It does not grant additional rights to the third-party AK135
publication or data.
