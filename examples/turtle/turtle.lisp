(define turtle-line-width 5)
(define turtle-pen-colour red)
(define turtle-pen-down #t)
(define turtle-pos (cons (/ win-w 2) (/ win-h 2)))
(define turtle-heading 0)

(define mod-fixed (lambda (n m)
    (if (> n 0) (mod n m) (mod (+ n m) m)))) 

(define turn (lambda (radians)
    (set! turtle-heading (mod-fixed (+ turtle-heading radians) (* 2 pi)))))
(define left (lambda () (turn (/ pi 2))))
(define right (lambda () (turn (- 0 (/ pi 2)))))

(define forward (lambda (n) (let ((px (car turtle-pos)) (py (cdr turtle-pos)) (vx (cos turtle-heading)) (vy (sin turtle-heading)))
    (let ((newx (+ px (* n vx))) (newy (+ py (* n vy))))
        (let ((newpos (cons newx newy)))
            (if (eqv? turtle-pen-down #t) (draw-line turtle-pos newpos))
            (set! turtle-pos newpos)
            (wraparound)
)))))

(define wraparound (lambda () (let ((px (car turtle-pos)) (py (cdr turtle-pos)))
    (set! turtle-pos (cons (mod-fixed px win-w) (mod-fixed py win-h)))
)))

(define draw-line (lambda (from to)
    (let ((imd (imdraw)))
        (im-set-color! imd turtle-pen-colour)
        (im-push imd (vec2d (car from) (cdr from)))
        (im-push imd (vec2d (car to) (cdr to)))
        (line imd turtle-line-width)
        (im-draw imd win)
        (update win)
)))

(define gosper3 (quote (f r f r r f l f l l f f l f r r l f r f f r r f r f l l f l f r r l f r f f r r f r f l l f l f l f r f r r f l f l l f f l f r l l f r f r r f l f l l f f l f r f r f r r f l f l l f f l f r l l f r f f r r f r f l l f l f r r l f r f r r f l f l l f f l f r r l f r f f r r f r f l l f l f l f r f f r r f r f l l f l f r r l f r f f r r f r f l l f l f r f r f r r f l f l l f f l f r l l f r f r r f l f l l f f l f r l l f r f f r r f r f l l f l f r r l f r f r r f l f l l f f l f r r l f r f f r r f r f l l f l f l f r f f r r f r f l l f l f r r l f r f f r r f r f l l f l f r f r f r r f l f l l f f l f r l l f r f r r f l f l l f f l f r l l f r f f r r f r f l l f l f l f r f r r f l f l l f f l f r r l f r f f r r f r f l l f l f r r l f r f f r r f r f l l f l f l f r f r r f l f l l f f l f r l l f r f r r f l f l l f f l f r f r f r r f l f l l f f l f r l l f r f f r r f r f l l f l f r l l f r f r r f l f l l f f l f r r l f r f f r r f r f l l f l f r r l f r f f r r f r f l l f l f l f r f r r f l f l l f f l f r l l f r f r r f l f l l f f l f r f r f r r f l f l l f f l f r l l f r f f r r f r f l l f l f r f r f r r f l f l l f f l f r r l f r f f r r f r f l l f l f r r l f r f f r r f r f l l f l f l f r f r r f l f l l f f l f r l l f r f r r f l f l l f f l f r f r f r r f l f l l f f l f r l l f r f f r r f r f l l f l f r l l f r f r r f l f l l f f l f r r l f r f f r r f r f l l f l f l f r f f r r f r f l l f l f r r l f r f f r r f r f l l f l f r f r f r r f l f l l f f l f r l l f r f r r f l f l l f f l f r l l f r f f r r f r f l l f l f r)))
(define program (quote ()))

(define start (lambda ()
    (begin
    (set! program gosper3)
    (set! tick (lambda () (begin
        (if (null? program) (set! program (quote (done))))
        (let ((next (car program)) (rem (cdr program)))
            (cond
                ((eqv? next 'f) (forward 10))
                ((eqv? next 'l)  (turn (/ pi 3)))
                ((eqv? next 'r) (turn (- 0 (/ pi 3))))
                (else #t))
            (set! program rem)
)))))))
