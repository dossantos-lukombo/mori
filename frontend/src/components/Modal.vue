<template>
    <div class="modal">
        <div class="modal-content">
            <div class="modal-header">
                <h2>
                    <slot name="title"></slot>
                </h2>
                <i class="uil uil-times close" @click="$emit('closeModal')"></i>
            </div>
            <div class="modal-body">
                <slot name="body"></slot>
            </div>
        </div>
    </div>
</template>


<style>
.modal {
    position: fixed;
    height: 100vh;
    width: 100vw;
    top: 0;
    left: 0;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    animation: fadeIn 0.3s ease-out;
}

.modal-content {
    display: flex;
    flex-direction: column;
    width: 90%;
    max-width: 450px;
    color: var(--color-white);
    position: relative;
    background-color: var(--bg-neutral);
    border-radius: 10px;
    overflow: hidden;
    box-shadow: 0 5px 20px rgba(0, 0, 0, 0.3);
    animation: slideUp 0.3s ease-out;
}

.modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.modal-header h2 {
    margin: 0;
    font-size: 1.3rem;
}

.modal-body {
    padding: 20px;
}

.modal-content form {
    display: flex;
    flex-direction: column;
    gap: 15px;
    width: 100%;
}

.modal-content .btn.form-submit {
    align-self: flex-end;
}

/* close icon */
.uil-times.close {
    cursor: pointer;
    font-size: 1.5rem;
    color: var(--color-white);
    opacity: 0.7;
    transition: opacity 0.2s ease;
    position: absolute;
    right: 20px;
    top: 20px;
}

.uil-times.close:hover {
    opacity: 1;
}

@keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
}

@keyframes slideUp {
    from { transform: translateY(30px); opacity: 0; }
    to { transform: translateY(0); opacity: 1; }
}

/* Add these new styles to ensure icons are properly centered */
.modal-actions .btn i,
button i {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    vertical-align: middle;
    width: 1.2em;
    height: 1.2em;
}

.modal button {
    display: inline-flex;
    align-items: center;
    justify-content: center;
}
</style>